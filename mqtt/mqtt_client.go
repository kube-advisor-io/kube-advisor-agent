package mqtt

import (
	"crypto/tls"
	"fmt"
	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
	"crypto/x509"
)

type MQTTClient struct {
	client          mqtt.Client
	qos             int
	previousMessage string
}

type MQTTOptions struct {
	clientOpts *mqtt.ClientOptions
	qos        int
}

func ParseConfig(config config.MQTTConfig) *MQTTOptions {
	fmt.Printf("MQTT Config:\n")
	fmt.Printf("\tbroker:       %s\n", config.Broker)
	fmt.Printf("\ttopic:        %s\n", config.Topic)
	fmt.Printf("\tusername:     %s\n", config.Username)
	fmt.Printf("\tpassword:     %s\n", config.Password)
	fmt.Printf("\tclientID:     %s\n", config.ClientID)
	fmt.Printf("\tkey:          %s\n", config.TlsKeyFile)
	fmt.Printf("\tcert:         %s\n", config.TlsCertificateFile)
	fmt.Printf("\tqos:          %d\n", config.Qos)
	fmt.Printf("\tcleanSession: %v\n", config.CleanSession)

	clientOpts := mqtt.NewClientOptions()
	clientOpts.AddBroker(config.Broker)
	if config.TlsKeyFile != "" && config.TlsCertificateFile != "" {
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(config.CACertificate)
		if err != nil {
			fmt.Println(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		cert, err := tls.LoadX509KeyPair(config.TlsCertificateFile, config.TlsKeyFile)
		if err != nil {
			fmt.Println(err.Error())
		}
		clientOpts.SetTLSConfig(&tls.Config{
			RootCAs: certpool,
			Certificates: []tls.Certificate{cert},
			ClientAuth: tls.NoClientCert,
			ClientCAs: nil,
		})
	}
	clientOpts.SetMaxReconnectInterval(1 * time.Second)
	clientOpts.SetKeepAlive(30 * time.Second)
	if config.Username != "" {
		fmt.Println("Set username "+ config.Username)
		clientOpts.SetUsername(config.Username)
		clientOpts.SetPassword(config.Password)
	}
	if config.ClientID != "" {
		clientOpts.SetClientID(config.ClientID)
		fmt.Println("Set client id "+ config.ClientID)
	}
	// opts.SetClientID(*id)

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
