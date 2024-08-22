package mqtt

import (
	"crypto/tls"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
	qos    int
	previousMessage string
}

func ParseMQTTFlags() (*mqtt.ClientOptions, int) {
	topic := flag.String("topic", "robert/robertstestsensor/message/testmessage", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://test.mosquitto.org:1884", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "readwrite", "The password (optional)")
	user := flag.String("user", "rw", "The User (optional)")
	clientKey := flag.String("clientKey", "", "The path to the key for the client certificate")
	clientCert := flag.String("clientCert", "", "The path to the client certificate")
	cleansess := flag.Bool("clean", true, "Set Clean Session (default false)")
	qos := *flag.Int("qos", 2, "The Quality of Service 0,1,2 (default 0)")
	flag.Parse()

	fmt.Printf("Sample Info:\n")
	fmt.Printf("\tbroker:    %s\n", *broker)
	fmt.Printf("\tuser:      %s\n", *user)
	fmt.Printf("\tpassword:  %s\n", *password)
	fmt.Printf("\tkey:       %s\n", *clientKey)
	fmt.Printf("\tcert:      %s\n", *clientCert)
	fmt.Printf("\ttopic:     %s\n", *topic)
	fmt.Printf("\tqos:       %d\n", qos)
	fmt.Printf("\tcleansess: %v\n", *cleansess)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(*broker)
	if *clientKey != "" && *clientCert != "" {
		cert, _ := tls.LoadX509KeyPair(*clientCert, *clientKey)
		opts.SetTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})
	}
	// opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)
	return opts, qos
}

func StartNewMQTTClient(opts *mqtt.ClientOptions, qos int) *MQTTClient {
	mqttClient := new(MQTTClient)
	mqttClient.client = mqtt.NewClient(opts)
	if token := mqttClient.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mqttClient.qos = qos
	fmt.Printf("MQTT Client Publisher started with opts %v, QOS %v\n", *opts, qos)
	return mqttClient
}

func (mqttClient *MQTTClient) PublishMessage(topic string, message string) {
	if (mqttClient.previousMessage == message){
		return
	}

	token := mqttClient.client.Publish(
		topic,
		byte(mqttClient.qos),
		false,
		message,
	)
	token.Wait()
	mqttClient.previousMessage = message
}
