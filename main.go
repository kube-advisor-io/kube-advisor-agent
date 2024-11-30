package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/mqtt"
	"github.com/bobthebuilderberlin/kube-advisor-agent/providers"
	"github.com/go-logr/logr"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	controllerlog "sigs.k8s.io/controller-runtime/pkg/log"
)

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		log.Error(err)
		return
	}

	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	controllerlog.SetLogger(logr.Discard())

	options, err := mqtt.ParseConfig(config.MQTT)
	if err != nil {
		log.Error(err)
		return
	}

	mqttClient, err := mqtt.StartNewMQTTClient(options)
	if err != nil {
		log.Error(err)
		return
	}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Error(err)
		return
	}

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
	kyvernoProvider := providers.NewKyvernoPoliciesProvider(dynamicClient, kubeConfig, config)

	gatherDataAndPublish(dataProviders, resourceProviders, kyvernoProvider, mqttClient, config)
	for range time.Tick(time.Second * 30) {
		gatherDataAndPublish(dataProviders, resourceProviders, kyvernoProvider, mqttClient, config)
	}

	waitIndefinitely()
}

func gatherDataAndPublish(
	dataProviders *[]DataProvider,
	resourceProviders *[]ResourceProvider,
	kyvernoProvider *providers.KyvernoPoliciesProvider,
	mqttClient *mqtt.MQTTClient,
	config config.Config,
) {
	messageData := make(map[string]interface{})
	messageData["id"] = config.OrganizationId + "_" + config.ClusterId
	messageData["version"] = 2 // schema version
	messageData["organizationId"] = config.OrganizationId
	messageData["clusterId"] = config.ClusterId

	data := make(map[string]interface{})

	for _, dataProvider := range *dataProviders {
		providerData := dataProvider.GetData()
		data[dataProvider.GetName()] = providerData
	}
	for _, resourceProvider := range *resourceProviders {
		parsedItems := resourceProvider.GetParsedItems()
		result := map[string]interface{}{}
		result["version"] = resourceProvider.GetVersion()
		result["items"] = parsedItems
		data[resourceProvider.GetResource().Resource] = result
	}

	clusterReport := kyvernoProvider.CheckPolicies()
	data["kyvernoClusterReports"] = clusterReport

	messageData["data"] = data

	jsonString, err := json.Marshal(messageData)
	if err != nil {
		log.Error(err)
		return
	}

	mqttClient.PublishMessage(config.MQTT.Topic, string(jsonString))
}

func waitIndefinitely() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
