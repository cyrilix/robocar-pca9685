package part

import (
	"encoding/json"
	"fmt"
	"github.com/cyrilix/robocar-base/service"
	"github.com/cyrilix/robocar-base/types"
	"github.com/cyrilix/robocar-pca9685/actuator"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

type Pca9685Part struct {
	client       MQTT.Client
	throttleCtrl *actuator.Throttle
	steeringCtrl *actuator.Steering

	muSteering    sync.Mutex
	steeringValue float64
	muThrottle    sync.Mutex
	throttleValue float64

	updateFrequency int

	throttleTopic string
	steeringTopic string

	cancel chan interface{}
}

func NewPca9685Part(client MQTT.Client, throttleCtrl *actuator.Throttle, steeringCtrl *actuator.Steering, updateFrequency int, throttleTopic, steeringTopic string) *Pca9685Part {
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
	ticker := time.NewTicker(time.Second / time.Duration(p.updateFrequency))

	for {
		select {
		case <-ticker.C:
			p.updateCtrl()
		case <-p.cancel:
			return nil
		}
	}
}

func (p *Pca9685Part) Stop() {
	close(p.cancel)
	service.StopService("pca9685", p.client, p.throttleTopic, p.steeringTopic)
}

func (p *Pca9685Part) onThrottleChange(_ MQTT.Client, message MQTT.Message) {
	var throttle types.Throttle
	err := json.Unmarshal(message.Payload(), &throttle)
	if err != nil {
		log.Printf("[%v] unable to unmarshall throttle msg: %v", message.Topic(), err)
		return
	}
	p.muThrottle.Lock()
	defer p.muThrottle.Unlock()
	p.throttleCtrl.SetPercentValue(throttle.Value)
}

func (p *Pca9685Part) onSteeringChange(_ MQTT.Client, message MQTT.Message) {
	var steering types.Steering
	err := json.Unmarshal(message.Payload(), &steering)
	if err != nil {
		log.Printf("[%v] unable to unmarshall steering msg: %v", message.Topic(), err)
		return
	}
	p.muSteering.Lock()
	defer p.muSteering.Unlock()
	p.steeringCtrl.SetPercentValue(steering.Value)
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

func (p *Pca9685Part) updateCtrl() {
	p.muThrottle.Lock()
	defer p.muThrottle.Unlock()
	p.throttleCtrl.SetPercentValue(p.throttleValue)

	p.muSteering.Lock()
	defer p.muSteering.Unlock()
	p.steeringCtrl.SetPercentValue(p.steeringValue)
}
