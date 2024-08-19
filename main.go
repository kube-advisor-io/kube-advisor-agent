package main

import (
	// "context"
	"encoding/json"
	"fmt"
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
	"k8s.io/client-go/tools/clientcmd"
	// toolsWatch "k8s.io/client-go/tools/watch"
)

var (
	config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ = kubernetes.NewForConfig(config)
)

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
	mqttClient := mqtt.StartNewMQTTClient(mqtt.ParseMQTTFlags())
	dataproviders := getAllDataProviders(clientset)

	for range time.Tick(time.Second * 10) {
		gatherDataAndPublish(dataproviders, mqttClient)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func gatherDataAndPublish(dataproviders []DataProvider, mqttClient *mqtt.MQTTClient) {
	data := make(map[string]interface{})
	for _, dataprovider := range dataproviders {
		maps.Copy(data, dataprovider.GetData())
	}
	jsonString, _ := json.Marshal(data)
	fmt.Println(jsonString)
	mqttClient.PublishMessage("robert/robertstestsensor/message/testmessage", jsonString)
	fmt.Println("Published data")
}
