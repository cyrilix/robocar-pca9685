package actuator

import (
	"fmt"
	"github.com/cyrilix/robocar-pca9685/pkg/util"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/pca9685"
	"periph.io/x/host/v3"
)

const (
	MinPercent = -1.
	MaxPercent = 1.
)

type PWM int

func NewDevice(freq physic.Frequency) *pca9685.Dev {
	zap.S().Info("NewDevice pca9685 controller")
	_, err := host.Init()
	if err != nil {
		zap.S().Fatalf("unable to NewDevice host: %v", err)
	}

	zap.S().Info("open i2c bus")
	bus, err := i2creg.Open("")
	if err != nil {
		zap.S().Fatalf("unable to NewDevice i2c bus: %v", err)
	}
	zap.S().Infof("i2c bus opened: %v", bus)

	device, err := pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		zap.S().Fatalf("unable to NewDevice pca9685 bus: %v", err)
	}
	zap.S().Infof("set pwm frequency to %d", 60)
	err = device.SetPwmFreq(freq)
	if err != nil {
		zap.S().Panicf("unable to set pwm frequency: %v", err)
	}
	zap.S().Info("NewDevice done")
	return device
}

func convertToDuty(percent float32, freq physic.Frequency, centerPWM, minPWM, maxPWM PWM) (gpio.Duty, error) {
	// map absolute angle to angle that vehicle can implement.
	pw := int(centerPWM)
	if percent > 0 {
		pw = util.MapRange(float64(percent), 0, MaxPercent, float64(centerPWM), float64(maxPWM))
	} else if percent < 0 {
		pw = util.MapRange(float64(percent), MinPercent, 0, float64(minPWM), float64(centerPWM))
	}
	zap.S().Debugf("convert value %v to pw: %v", percent, pw)

	per := freq.Period().Microseconds()

	draw := float64(pw) / float64(per)
	d, err := gpio.ParseDuty(fmt.Sprintf("%.f%%", draw*100))
	if err != nil {
		return 0, fmt.Errorf("unable to parse duty, probably bad compute: %w", err)
	}
	return d, nil
}

type Pca9685Controller struct {
	pin                        gpio.PinIO
	minPWM, maxPWM, neutralPWM PWM
	freq                       physic.Frequency
}

func (c *Pca9685Controller) Close() error {
	return c.pin.Halt()
}

func (c *Pca9685Controller) SetDuty(d gpio.Duty) error {
	err := c.pin.PWM(d, c.freq)
	if err != nil {
		return fmt.Errorf("unable to set pwm value: %v", err)
	}
	return nil
}

// SetPercentValue Set percent value
func (c *Pca9685Controller) SetPercentValue(p float32) error {
	d, err := convertToDuty(p, c.freq, c.neutralPWM, c.minPWM, c.maxPWM)
	if err != nil {
		return fmt.Errorf("unable to compute Duty value for steering: %w", err)
	}
	err = c.SetDuty(d)
	if err != nil {
		return fmt.Errorf("unable to apply duty value '%v' for pin '%v': '%w' ", d, c.pin.Name(), err)
	}
	return nil
}

func NewPca9685Controller(device *pca9685.Dev, channel int, minPWM, maxPWM, neutralPWM PWM, freq physic.Frequency) (*Pca9685Controller, error) {
	p, err := device.CreatePin(channel)
	if err != nil {
		return nil, fmt.Errorf("unable to create pin for channel %v: %w", channel, err)
	}
	s := Pca9685Controller{
		pin:        p,
		minPWM:     minPWM,
		maxPWM:     maxPWM,
		neutralPWM: neutralPWM,
		freq:       freq,
	}
	return &s, nil
}
