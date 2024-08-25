package mqtt

import (
	"crypto/tls"
	"fmt"
	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client          mqtt.Client
	qos             int
	previousMessage string
}

type MQTTOptions struct {
	clientOpts *mqtt.ClientOptions
	qos int
}

func ParseConfig(config config.MQTTConfig) *MQTTOptions {
	fmt.Printf("MQTT Config:\n")
	fmt.Printf("\tbroker:       %s\n", config.Broker)
	fmt.Printf("\topic:         %s\n", config.Topic)
	fmt.Printf("\tusername:     %s\n", config.Username)
	fmt.Printf("\tpassword:     %s\n", config.Password)
	fmt.Printf("\tkey:          %s\n", config.TlsKeyFile)
	fmt.Printf("\tcert:         %s\n", config.TlsCertificateFile)
	fmt.Printf("\tqos:          %d\n", config.Qos)
	fmt.Printf("\tcleanSession: %v\n", config.CleanSession)

	clientOpts := mqtt.NewClientOptions()
	clientOpts.AddBroker(config.Broker)
	if config.TlsKeyFile != "" && config.TlsCertificateFile != "" {
		cert, _ := tls.LoadX509KeyPair(config.TlsKeyFile, config.TlsCertificateFile)
		clientOpts.SetTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})
	}
	// opts.SetClientID(*id)
	clientOpts.SetUsername(config.Username)
	clientOpts.SetPassword(config.Password)
	clientOpts.SetCleanSession(config.CleanSession)
	return &MQTTOptions{clientOpts: clientOpts, qos: config.Qos}
}

func StartNewMQTTClient(opts *MQTTOptions) *MQTTClient {
	mqttClient := new(MQTTClient)
	mqttClient.client = mqtt.NewClient(opts.clientOpts)
	if token := mqttClient.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mqttClient.qos = opts.qos
	fmt.Printf("MQTT Client Publisher started with opts %v", *opts.clientOpts)
	return mqttClient
}

func (mqttClient *MQTTClient) PublishMessage(topic string, message string) {
	fmt.Printf("Trying to publish data %v ...", message)
	if mqttClient.previousMessage == message {
		fmt.Println("was already sent")
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
	fmt.Println("published.")
}
