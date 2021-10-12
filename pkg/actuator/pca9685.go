package actuator

import (
	"go.uber.org/zap"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/pca9685"
	"periph.io/x/host/v3"
)

var (
	device *pca9685.Dev
)

func init() {
	zap.S().Info("init pca9685 controller")
	_, err := host.Init()
	if err != nil {
		zap.S().Fatalf("unable to init host: %v", err)
	}

	zap.S().Info("open i2c bus")
	bus, err := i2creg.Open("")
	if err != nil {
		zap.S().Fatalf("unable to init i2c bus: %v", err)
	}
	zap.S().Info("i2c bus opened")

	device, err = pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		zap.S().Fatalf("unable to init pca9685 bus: %v", err)
	}
	zap.S().Infof("set pwm frequency to %d", 60)
	err = device.SetPwmFreq(60 * physic.Hertz)
	if err != nil {
		zap.S().Panicf("unable to set pwm frequency: %v", err)
	}
	zap.S().Info("init done")
}
