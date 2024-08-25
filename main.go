package main

import (
	// "context"
	"encoding/json"
	// "fmt"
	"maps"
	"os"
	"sync"
	"time"
	// log "github.com/sirupsen/logrus"
	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/watch"
	"github.com/bobthebuilderberlin/kube-advisor-agent/mqtt"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/cache"
	config "github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"k8s.io/client-go/tools/clientcmd"
	// toolsWatch "k8s.io/client-go/tools/watch"
)

var ()

// func watchNamespaces() {

//     watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
//         timeOut := int64(60)
// 		log.Info("Starting watching namespaces")
//         return clientset.CoreV1().Namespaces().Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut})
//     }

//     watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})

//     for event := range watcher.ResultChan() {
//         item := event.Object.(*corev1.Namespace)

//         switch event.Type {
//         case watch.Modified:
//         case watch.Bookmark:
//         case watch.Error:
//         case watch.Deleted:
//         case watch.Added:
//             processNamespace(item.GetName())
//         }
//     }
// }

// func processNamespace(namespace string) {
//     log.Info("Some processing for newly created namespace : ", namespace)
// }

// func watchPods() {

//     watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
//         timeOut := int64(60)
// 		log.Info("Starting watching namespaces")
//         return clientset.CoreV1().Pods("").Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut})
//     }

//     watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})

//     for event := range watcher.ResultChan() {
//         item := event.Object.(*corev1.Pod)

//         switch event.Type {
//         case watch.Modified:
//         case watch.Bookmark:
//         case watch.Error:
//         case watch.Deleted:
//         case watch.Added:
//             processPod(item.GetName())
//         }
//     }
// }

// func processPod(pod string) {
//     log.Info("Some processing for newly created pod: ", pod)
//     token := client.Publish("robertssupercoolclientidthatisnotused", 2, false, pod)
//     token.Wait()
// }

func main() {
	config := config.ReadConfig()
	mqttClient := mqtt.StartNewMQTTClient(mqtt.ParseConfig(config.MQTT))

	kubeConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ := kubernetes.NewForConfig(kubeConfig)
	dataproviders := getAllDataProviders(clientset, config.DisabledProviders)

	for range time.Tick(time.Second * 10) {
		gatherDataAndPublish(dataproviders, mqttClient, config)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func gatherDataAndPublish(dataproviders *[]DataProvider, mqttClient *mqtt.MQTTClient, config config.Config) {
	data := make(map[string]interface{})
	for _, dataprovider := range *dataproviders {
		maps.Copy(data, dataprovider.GetData())
	}
	data["id"] = config.CustomerID + "_" + config.ClusterID
	data["version"] = "0.1" // schema version
	jsonString, _ := json.Marshal(data)
	mqttClient.PublishMessage(config.MQTT.Topic, string(jsonString))
}
