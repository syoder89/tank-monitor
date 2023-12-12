package main

import (
	"fmt"
	"time"
	"os"
	"encoding/json"
	"github.com/VictoriaMetrics/metrics"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/syoder89/tank-monitor/vmclient"
)

type TankMsg struct {
	Distance float64
	Temperature float64
	Humidity float64
}

var tmsg TankMsg
var sensor string
// tcp://mosquitto
var broker = "tcp://mosquitto:1883"
// http://172.20.1.4:8428/api/v1/import/prometheus
var vmPushURL = "http://victoria-metrics-victoria-metrics-single-server:8428/api/v1/import/prometheus"

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	json.Unmarshal([]byte(msg.Payload()), &tmsg)
	fmt.Println(tmsg)
	vmclient.Push(vmPushURL, 20*time.Second, `sensor="`+sensor+`"`, false)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if val, ok := os.LookupEnv("SENSOR"); ok {
		sensor = val
	} else {
		panic("No sensor name provided!")
	}

	if val, ok := os.LookupEnv("BROKER"); ok {
		broker = val
	}
	if val, ok := os.LookupEnv("VM_PUSH_URL"); ok {
		vmPushURL = val
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("monitor-"+sensor)
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetCleanSession(true)
	opts.SetOrderMatters(false)
	opts.KeepAlive(30)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	metrics.NewGauge(`distance`, func() float64 { return tmsg.Distance })
	metrics.NewGauge(`temperature`, func() float64 { return tmsg.Temperature })
	metrics.NewGauge(`humidity`, func() float64 { return tmsg.Humidity })

	topic := "tele/"+sensor+"/SENSOR"
	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *broker)
	}

	<-c
}

// Received message: {"Distance": 1000,"Temperature": 23.760967,"Humidity": 33.665981} from topic: tele/taylor_water_tank_level1/SENSOR

