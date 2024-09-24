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
	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	mqttClient := mqtt.StartNewMQTTClient(mqtt.ParseConfig(config.MQTT))
	kubeConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	staticClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Error(err)
		return
	}
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		log.Error(err)
		return
	}
	dataProviders := getAllDataProviders(staticClient, config)
	resourceProviders := getAllResourceProviders(dynamicClient, config)
	for range time.Tick(time.Second * 10) {
		gatherDataAndPublish(dataProviders, resourceProviders, mqttClient, config)
	}

	waitIndefinitely()
}

func gatherDataAndPublish(dataProviders *[]DataProvider, resourceProviders *[]ResourceProvider, mqttClient *mqtt.MQTTClient, config config.Config) {
	messageData := make(map[string]interface{})
	messageData["id"] = config.OrganizationId + "_" + config.ClusterId
	messageData["version"] = 2 // schema version
	messageData["organizationId"] = config.OrganizationId
	messageData["clusterId"] = config.ClusterId

	data := make(map[string]interface{})

	for _, dataProvider := range *dataProviders {
		providerData := dataProvider.GetData()
		providerData["version"] = dataProvider.GetVersion()
		data[dataProvider.GetName()] = providerData
	}
	for _, resourceProvider := range *resourceProviders {
		parsedItems := resourceProvider.GetParsedItems()
		result := map[string]interface{}{}
		result["version"] = resourceProvider.GetVersion()
		result["items"] = parsedItems
		data[resourceProvider.GetResource().Resource] = result
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
