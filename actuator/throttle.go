package actuator

import (
	"github.com/cyrilix/robocar-pca9685/util"
	"log"
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
		log.Printf("unable to set throttle pwm value: %v", err)
	}

}

// Set percent value throttle
func (t *Throttle) SetPercentValue(p float64) {
	var pulse int
	if p > 0 {
		pulse = util.MapRange(p, 0, MaxThrottle, float64(t.zeroPulse), float64(t.maxPulse))
	} else {
		pulse = util.MapRange(p, MinThrottle, 0, float64(t.minPulse), float64(t.zeroPulse))
	}
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

	log.Printf("send zero pulse to calibrate ESC: %v", zeroPulse)
	t.SetPulse(zeroPulse)

	return &t
}
