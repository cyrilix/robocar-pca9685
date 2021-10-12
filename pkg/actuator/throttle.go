package actuator

import (
	"github.com/cyrilix/robocar-pca9685/pkg/util"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/devices/v3/pca9685"
)

const (
	MinThrottle = -1
	MaxThrottle = 1
)

type Throttle struct {
	channel                       int
	zeroPulse, minPulse, maxPulse int
	dev                           *pca9685.Dev
}

func (t *Throttle) SetPulse(pulse int) {
	err := t.dev.SetPwm(t.channel, 0, gpio.Duty(pulse))
	if err != nil {
		zap.S().Errorf("unable to set throttle pwm value: %v", err)
	}

}

// Set percent value throttle
func (t *Throttle) SetPercentValue(p float32) {
	var pulse int
	if p > 0 {
		pulse = util.MapRange(float64(p), 0, MaxThrottle, float64(t.zeroPulse), float64(t.maxPulse))
	} else {
		pulse = util.MapRange(float64(p), MinThrottle, 0, float64(t.minPulse), float64(t.zeroPulse))
	}
	zap.S().Debugf("set throttle to %v-> %v (%v, %v, %v, %v, %v)", p, pulse, LeftAngle, RightAngle, t.minPulse, t.maxPulse, t.zeroPulse)
	t.SetPulse(pulse)
}

func NewThrottle(channel, zeroPulse, minPulse, maxPulse int) *Throttle {
	t := Throttle{
		channel:   channel,
		dev:       device,
		zeroPulse: zeroPulse,
		minPulse:  minPulse,
		maxPulse:  maxPulse,
	}

	zap.S().Info("send zero pulse to calibrate ESC")
	t.SetPercentValue(0)

	return &t
}
