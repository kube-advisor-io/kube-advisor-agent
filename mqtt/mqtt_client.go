package mqtt

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
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

func ParseConfig(config config.MQTTConfig) (*MQTTOptions, error) {
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
		log.Info("Using TLS client certificate auth for MQTT connection.")
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(config.CACertificate)
		if err != nil {
			return nil, err
		}
		certpool.AppendCertsFromPEM(ca)
		cert, err := tls.LoadX509KeyPair(config.TlsCertificateFile, config.TlsKeyFile)
		if err != nil {
			return nil, err
		}
		clientOpts.SetTLSConfig(&tls.Config{
			RootCAs:      certpool,
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
			ClientCAs:    nil,
		})
	}
	clientOpts.SetMaxReconnectInterval(1 * time.Second)
	clientOpts.SetKeepAlive(30 * time.Second)
	if config.Username != "" {
		log.Info("Using username/password for MQTT connection. Username: " + config.Username)
		clientOpts.SetUsername(config.Username)
		clientOpts.SetPassword(config.Password)
	}
	if config.ClientID != "" {
		clientOpts.SetClientID(config.ClientID)
		log.Info("Set client id for MQTT connection: " + config.ClientID)
	}

	clientOpts.SetCleanSession(config.CleanSession)
	return &MQTTOptions{clientOpts: clientOpts, qos: config.Qos}, nil
}

func StartNewMQTTClient(opts *MQTTOptions) (*MQTTClient, error) {
	mqttClient := new(MQTTClient)
	mqttClient.client = mqtt.NewClient(opts.clientOpts)
	if token := mqttClient.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	mqttClient.qos = opts.qos
	log.Info("MQTT Client Publisher started with opts ", *opts.clientOpts)
	return mqttClient, nil
}

func (mqttClient *MQTTClient) PublishMessage(topic string, message string) {
	log.Info("Trying to publish data via MQTT ...")
	if mqttClient.previousMessage == message {
		log.Info("... but exact same data was already sent previously.")
		return
	}

	gzippedMessage, err := gzipMessage(message)
	if err != nil {
		log.Error("error gzipping message: ", err)
		return
	}

	token := mqttClient.client.Publish(
		topic,
		byte(mqttClient.qos),
		false,
		gzippedMessage,
	)
	token.Wait()

	if token.Error() != nil {
		log.Error("error publishing message to MQTT Broker: ", token.Error())
		return
	}

	mqttClient.previousMessage = message
	log.Info("Published message: ", message)
	log.Info(fmt.Sprintf("Message length: %v bytes, gzipped %v bytes.", len(message), len(gzippedMessage)))
}

func gzipMessage(message string) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(message)); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
