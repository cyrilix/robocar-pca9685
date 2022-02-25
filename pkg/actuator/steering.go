package actuator

import (
	"github.com/cyrilix/robocar-pca9685/pkg/util"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/devices/v3/pca9685"
)

const (
	LeftAngle  = -1.
	RightAngle = 1.
)

type Steering struct {
	channel                      int
	leftPWM, rightPWM, centerPWM int
	dev                          *pca9685.Dev
}

func (s *Steering) SetPulse(pulse int) {
	err := s.dev.SetPwm(s.channel, 0, gpio.Duty(pulse))
	if err != nil {
		zap.S().Warnf("unable to set steering pwm value: %v", err)
	}

}

// SetPercentValue Set percent value steering
func (s *Steering) SetPercentValue(p float32) {
	// map absolute angle to angle that vehicle can implement.

	pulse := s.centerPWM
	if p > 0 {
		pulse = util.MapRange(float64(p), 0, RightAngle, float64(s.centerPWM), float64(s.rightPWM))
	} else if p < 0 {
		pulse = util.MapRange(float64(p), LeftAngle, 0, float64(s.leftPWM), float64(s.centerPWM))
	}
	zap.S().Debugf("convert steering %v to %v", p, pulse)
	s.SetPulse(pulse)
}

func NewSteering(channel, leftPWM, rightPWM, centerPWM int) *Steering {
	s := Steering{
		channel:   channel,
		dev:       device,
		leftPWM:   leftPWM,
		rightPWM:  rightPWM,
		centerPWM: centerPWM,
	}
	return &s
}
