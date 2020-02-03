package actuator

import (
	"github.com/cyrilix/robocar-pca9685/util"
	log "github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/experimental/devices/pca9685"
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
		log.Infof("unable to set throttle pwm value: %v", err)
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
	log.Debugf("set throttle to %v-> %v (%v, %v, %v, %v, %v)", p, pulse, LeftAngle, RightAngle, t.minPulse, t.maxPulse, t.zeroPulse)
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

	log.Infof("send zero pulse to calibrate ESC")
	t.SetPulse(zeroPulse)

	return &t
}
