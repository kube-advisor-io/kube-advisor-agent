package main

import (
	"encoding/json"
	"maps"
	"os"
	"sync"
	"time"

	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/mqtt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config := config.ReadConfig()
	mqttClient := mqtt.StartNewMQTTClient(mqtt.ParseConfig(config.MQTT))

	kubeConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ := kubernetes.NewForConfig(kubeConfig)
	dataproviders := getAllDataProviders(clientset, config.DisabledProviders)

	for range time.Tick(time.Second * 10) {
		gatherDataAndPublish(dataproviders, mqttClient, config)
	}

	waitIndefinitely()
}

func gatherDataAndPublish(dataproviders *[]DataProvider, mqttClient *mqtt.MQTTClient, config config.Config) {
	data := make(map[string]interface{})
	for _, dataprovider := range *dataproviders {
		maps.Copy(data, dataprovider.GetData())
	}
	data["id"] = config.OrganizationID + "_" + config.ClusterID
	data["version"] = "1" // schema version
	data["organizationID"] = config.OrganizationID
	data["clusterID"] = config.ClusterID
	jsonString, _ := json.Marshal(data)
	mqttClient.PublishMessage(config.MQTT.Topic, string(jsonString))
}

func waitIndefinitely() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
