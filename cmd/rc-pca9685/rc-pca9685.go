package main

import (
	"flag"
	"github.com/cyrilix/robocar-base/cli"
	"github.com/cyrilix/robocar-pca9685/pkg/actuator"
	"github.com/cyrilix/robocar-pca9685/pkg/part"
	"go.uber.org/zap"
	"log"
	"os"
	"periph.io/x/conn/v3/physic"
)

const (
	DefaultClientId = "robocar-pca9685"
	SteeringChannel = 0
	ThrottleChannel = 1

	ThrottleStoppedPWM = 1455
	ThrottleMinPWM     = 1113
	ThrottleMaxPWM     = 1800

	SteeringLeftPWM  = 1004
	SteeringRightPWM = 1986
)

var (
	SteeringCenterPWM = (SteeringRightPWM-SteeringLeftPWM)/2 + SteeringLeftPWM
)

func main() {
	var mqttBroker, username, password, clientId, topicThrottle, topicSteering string

	mqttQos := cli.InitIntFlag("MQTT_QOS", 0)
	_, mqttRetain := os.LookupEnv("MQTT_RETAIN")

	cli.InitMqttFlags(DefaultClientId, &mqttBroker, &username, &password, &clientId, &mqttQos, &mqttRetain)

	var throttleChannel, throttleStoppedPWM, throttleMinPWM, throttleMaxPWM int
	if err := cli.SetIntDefaultValueFromEnv(&throttleChannel, "THROTTLE_CHANNEL", ThrottleChannel); err != nil {
		zap.S().Warnf("unable to init throttleChannel arg: %v", err)
	}
	if err := cli.SetIntDefaultValueFromEnv(&throttleStoppedPWM, "THROTTLE_STOPPED_PWM", ThrottleStoppedPWM); err != nil {
		zap.S().Warnf("unable to init throttleStoppedPWM arg: %v", err)
	}
	if err := cli.SetIntDefaultValueFromEnv(&throttleMinPWM, "THROTTLE_MIN_PWM", ThrottleMinPWM); err != nil {
		zap.S().Warnf("unable to init throttleMinPWM arg: %v", err)
	}
	if err := cli.SetIntDefaultValueFromEnv(&throttleMaxPWM, "THROTTLE_MAX_PWM", ThrottleMaxPWM); err != nil {
		zap.S().Warnf("unable to init throttleMaxPWM arg: %v", err)
	}

	var steeringChannel, steeringLeftPWM, steeringRightPWM, steeringCenterPWM int
	if err := cli.SetIntDefaultValueFromEnv(&steeringChannel, "STEERING_CHANNEL", SteeringChannel); err != nil {
		zap.S().Warnf("unable to init steeringChannel arg: %v", err)
	}
	if err := cli.SetIntDefaultValueFromEnv(&steeringLeftPWM, "STEERING_LEFT_PWM", SteeringLeftPWM); err != nil {
		zap.S().Warnf("unable to init steeringLeftPWM arg: %v", err)
	}
	if err := cli.SetIntDefaultValueFromEnv(&steeringRightPWM, "STEERING_RIGHT_PWM", SteeringRightPWM); err != nil {
		zap.S().Warnf("unable to init steeringRightPWM arg: %v", err)
	}
	if err := cli.SetIntDefaultValueFromEnv(&steeringCenterPWM, "STEERING_CENTER_PWM", SteeringCenterPWM); err != nil {
		zap.S().Warnf("unable to init steeringRightPWM arg: %v", err)
	}

	var updatePWMFrequency, pwmFreq int
	if err := cli.SetIntDefaultValueFromEnv(&updatePWMFrequency, "UPDATE_PWM_FREQUENCY", 25); err != nil {
		zap.S().Warnf("unable to init updatePWMFrequency arg: %v", err)
	}

	flag.StringVar(&topicThrottle, "mqtt-topic-throttle", os.Getenv("MQTT_TOPIC_THROTTLE"), "Mqtt topic that contains throttle value, use MQTT_TOPIC_THROTTLE if args not set")
	flag.StringVar(&topicSteering, "mqtt-topic-steering", os.Getenv("MQTT_TOPIC_STEERING"), "Mqtt topic that contains steering value, use MQTT_TOPIC_STEERING if args not set")
	flag.IntVar(&throttleChannel, "throttle-channel", throttleChannel, "I2C channel to use to control throttle, THROTTLE_CHANNEL env if args not set")
	flag.IntVar(&steeringChannel, "steering-channel", steeringChannel, "I2C channel to use to control steering, STEERING_CHANNEL env if args not set")
	flag.IntVar(&throttleStoppedPWM, "throttle-zero-pwm", throttleStoppedPWM, "Zero value for throttle PWM, THROTTLE_STOPPED_PWM env if args not set")
	flag.IntVar(&throttleMinPWM, "throttle-min-pwm", throttleMinPWM, "Left value for throttle PWM, THROTTLE_MIN_PWM env if args not set")
	flag.IntVar(&throttleMaxPWM, "throttle-max-pwm", throttleMaxPWM, "Right value for throttle PWM, THROTTLE_MAX_PWM env if args not set")
	flag.IntVar(&steeringLeftPWM, "steering-left-pwm", steeringLeftPWM, "Right left value for steering PWM, STEERING_LEFT_PWM env if args not set")
	flag.IntVar(&steeringRightPWM, "steering-right-pwm", steeringRightPWM, "Right right value for steering PWM, STEERING_RIGHT_PWM env if args not set")
	flag.IntVar(&steeringCenterPWM, "steering-center-pwm", steeringCenterPWM, "Center value for steering PWM, STEERING_CENTER_PWM env if args not set")
	flag.IntVar(&pwmFreq, "pwm-freq", 60, "PWM frequency in Hz")
	flag.IntVar(&updatePWMFrequency, "update-pwm-frequency", updatePWMFrequency, "Number of update values per seconds, UPDATE_PWM_FREQUENCY env if args not set")

	logLevel := zap.LevelFlag("log", zap.InfoLevel, "log level")

	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(*logLevel)
	lgr, err := config.Build()
	if err != nil {
		log.Fatalf("unable to init logger: %v", err)
	}
	defer func() {
		if err := lgr.Sync(); err != nil {
			log.Printf("unable to Sync logger: %v\n", err)
		}
	}()
	zap.ReplaceGlobals(lgr)
	client, err := cli.Connect(mqttBroker, username, password, clientId)
	if err != nil {
		zap.S().Fatalf("unable to connect to mqtt bus: %v", err)
	}
	defer client.Disconnect(50)

	freq := physic.Frequency(pwmFreq) * physic.Hertz

	zap.S().Infof("throttle channel  : %v", throttleChannel)
	zap.S().Infof("throttle frequency: %v", freq)
	zap.S().Infof("throttle zero     : %v", throttleStoppedPWM)
	zap.S().Infof("throttle min      : %v", throttleMinPWM)
	zap.S().Infof("throttle max      : %v", throttleMaxPWM)

	zap.S().Infof("steering channel  : %v", steeringChannel)
	zap.S().Infof("steering frequency: %v", freq)
	zap.S().Infof("steering center   : %v", steeringCenterPWM)
	zap.S().Infof("steering left     : %v", steeringLeftPWM)
	zap.S().Infof("steering right    : %v", steeringRightPWM)

	dev := actuator.NewDevice(freq)

	t, err := actuator.NewPca9685Controller(
		dev,
		throttleChannel,
		actuator.PWM(throttleMinPWM), actuator.PWM(throttleMaxPWM), actuator.PWM(throttleStoppedPWM),
		freq,
	)
	if err != nil {
		zap.S().Panicf("unable to init throttle controller: %v", err)
	}

	s, err := actuator.NewPca9685Controller(
		dev,
		steeringChannel,
		actuator.PWM(steeringLeftPWM), actuator.PWM(steeringRightPWM), actuator.PWM(steeringCenterPWM),
		freq,
	)
	if err != nil {
		zap.S().Panicf("unable to init steering controller: %v", err)
	}

	p := part.NewPca9685Part(client, t, s, updatePWMFrequency, topicThrottle, topicSteering)

	cli.HandleExit(p)

	zap.S().Info("devices ready, start event listener")
	err = p.Start()
	if err != nil {
		zap.S().Fatalf("unable to start service: %v", err)
	}
}
