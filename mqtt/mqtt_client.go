package mqtt

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
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
	log.Info("MQTT Config:\n")
	log.Info("broker:       ", config.Broker)
	log.Info("topic:        ", config.Topic)
	log.Info("username:     ", config.Username)
	log.Info("password:     ", config.Password)
	log.Info("clientID:     ", config.ClientID)
	log.Info("key:          ", config.TlsKeyFile)
	log.Info("cert:         ", config.TlsCertificateFile)
	log.Info("qos:          ", config.Qos)
	log.Info("cleanSession: ", config.CleanSession)

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
		log.Info("Set client id "+ config.ClientID)
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
	log.Info("MQTT Client Publisher started with opts ", *opts.clientOpts)
	return mqttClient
}

func (mqttClient *MQTTClient) PublishMessage(topic string, message string) {
	log.Info(fmt.Sprintf("Trying to publish data %v ...", message))
	if mqttClient.previousMessage == message {
		log.Info("was already sent")
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
	log.Info("published.")
}
