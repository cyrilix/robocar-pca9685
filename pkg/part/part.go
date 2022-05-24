package part

import (
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-protobuf/go/events"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
)

type ActuatorController interface {
	// SetPercentValue Set percent value
	SetPercentValue(p float32) error
	io.Closer
}

type Pca9685Part struct {
	client       MQTT.Client
	throttleCtrl ActuatorController
	steeringCtrl ActuatorController

	muSteering    sync.Mutex
	steeringValue float32
	muThrottle    sync.Mutex
	throttleValue float32

	updateFrequency int

	throttleTopic string
	steeringTopic string

	cancel chan interface{}
}

func NewPca9685Part(client MQTT.Client, throttleCtrl ActuatorController, steeringCtrl ActuatorController, updateFrequency int, throttleTopic, steeringTopic string) *Pca9685Part {
	return &Pca9685Part{
		client:          client,
		throttleCtrl:    throttleCtrl,
		steeringCtrl:    steeringCtrl,
		updateFrequency: updateFrequency,
		throttleTopic:   throttleTopic,
		steeringTopic:   steeringTopic,
		cancel:          make(chan interface{}),
	}
}

func (p *Pca9685Part) Start() error {
	if err := p.registerCallbacks(); err != nil {
		return fmt.Errorf("unable to start service: %v", err)
	}

	if err := p.steeringCtrl.SetPercentValue(0); err != nil {
		return fmt.Errorf("unable to start steering controller: %w", err)
	}
	if err := p.throttleCtrl.SetPercentValue(0); err != nil {
		return fmt.Errorf("unable to set init value '0%%': %w", err)
	}

	for {
		select {
		case <-p.cancel:
			return nil
		}
	}
}

func (p *Pca9685Part) Stop() {
	close(p.cancel)
	if err := p.steeringCtrl.Close(); err != nil {
		zap.S().Errorf("unable to close steering controller: %v", err)
	}
	if err := p.throttleCtrl.Close(); err != nil {
		zap.S().Errorf("unable to close throttle controller: %v", err)
	}
	service.StopService("pca9685", p.client, p.throttleTopic, p.steeringTopic)
}

func (p *Pca9685Part) onThrottleChange(_ MQTT.Client, message MQTT.Message) {
	var throttle events.ThrottleMessage
	err := proto.Unmarshal(message.Payload(), &throttle)
	if err != nil {
		zap.S().Warnw("unable to unmarshall throttle msg", "topic",
			message.Topic(),
			"error", err)
		return
	}
	zap.S().Debugf("new throttle value: %v", throttle.GetThrottle())
	p.muThrottle.Lock()
	defer p.muThrottle.Unlock()
	err = p.throttleCtrl.SetPercentValue(throttle.GetThrottle())
	if err != nil {
		zap.S().Errorf("unable to set steering value '%v': %w", throttle.GetThrottle(), err)
	}
}

func (p *Pca9685Part) onSteeringChange(_ MQTT.Client, message MQTT.Message) {
	var steering events.SteeringMessage
	err := proto.Unmarshal(message.Payload(), &steering)
	if err != nil {
		zap.S().Warnw("unable to unmarshal steering msg",
			"topic", message.Topic(),
			"error", err,
		)
		return
	}
	zap.S().Debugf("new steering value: %v", steering.GetSteering())
	p.muSteering.Lock()
	defer p.muSteering.Unlock()
	err = p.steeringCtrl.SetPercentValue(steering.GetSteering())
	if err != nil {
		zap.S().Errorf("unable to set steering value '%v': %w", steering.GetSteering(), err)
	}
}

func (p *Pca9685Part) registerCallbacks() error {
	err := service.RegisterCallback(p.client, p.throttleTopic, p.onThrottleChange)
	if err != nil {
		return fmt.Errorf("unable to register throttle callback: %v", err)
	}

	err = service.RegisterCallback(p.client, p.steeringTopic, p.onSteeringChange)
	if err != nil {
		return fmt.Errorf("unable to register steering callback: %v", err)
	}

	return nil
}
