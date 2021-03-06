package actuator

import (
	log "github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/pca9685"
	"periph.io/x/periph/host"
)

var (
	device *pca9685.Dev
)

func init() {
	log.Info("init pca9685 controller")
	_, err := host.Init()
	if err != nil {
		log.Fatalf("unable to init host: %v", err)
	}

	log.Info("open i2c bus")
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("unable to init i2c bus: %v", err)
	}
	log.Info("i2c bus opened")

	device, err = pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		log.Fatalf("unable to init pca9685 bus: %v", err)
	}
	log.Infof("set pwm frequency to %d", 60)
	err = device.SetPwmFreq(60 * physic.Hertz)
	if err != nil {
		log.Panicf("unable to set pwm frequency: %v", err)
	}
	log.Info("init done")
}
