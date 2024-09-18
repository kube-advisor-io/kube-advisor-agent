package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/mqtt"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config := config.ReadConfig()
	logLevel, err :=  log.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	mqttClient := mqtt.StartNewMQTTClient(mqtt.ParseConfig(config.MQTT))
	kubeConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ := kubernetes.NewForConfig(kubeConfig)
	dynamicClient, _ := dynamic.NewForConfig(kubeConfig)
	dataproviders := getAllDataProviders(clientset, dynamicClient, config)

	for range time.Tick(time.Second * 10) {
		gatherDataAndPublish(dataproviders, mqttClient, config)
	}

	waitIndefinitely()
}

func gatherDataAndPublish(dataproviders *[]DataProvider, mqttClient *mqtt.MQTTClient, config config.Config) {
	messageData := make(map[string]interface{})
	messageData["id"] = config.OrganizationID + "_" + config.ClusterID
	messageData["version"] = 2 // schema version
	messageData["organizationID"] = config.OrganizationID
	messageData["clusterID"] = config.ClusterID
	
	data := make(map[string]interface{})
	for _, dataprovider := range *dataproviders {
		providerData := dataprovider.GetData()
		providerData["version"] = dataprovider.GetVersion()
		data[dataprovider.GetName()] = providerData
	}
	messageData["data"] = data
	jsonString, _ := json.Marshal(messageData)
	mqttClient.PublishMessage(config.MQTT.Topic, string(jsonString))
}

func waitIndefinitely() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
