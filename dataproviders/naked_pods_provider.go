package dataproviders

import (
	"context"
	"time"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)
var (
	listTimeout = int64(60)
)

type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	owner     string
}

type NakedPodsProvider struct {
	client *kubernetes.Clientset
	pods   []*PodInfo
}

func NewNakedPodsProvider(client *kubernetes.Clientset) *NakedPodsProvider {
	instance := new(NakedPodsProvider)
	instance.client = client
	err := instance.updatePods()
	if err != nil {
		log.Error("Error updating pod list: ", err)
	}
	instance.startWatching()
	return instance
}

func (npp *NakedPodsProvider) startWatching() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go npp.tick(ticker, quit)
}

func (npp *NakedPodsProvider) tick(ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
			case <- ticker.C:
				err := npp.updatePods()
				if err != nil {
					log.Error("Error updating pod list: ", err)
				}
			case <- quit:
				ticker.Stop()
				return
			}
	}
}

func (npp *NakedPodsProvider) updatePods() error{
	pods, err := npp.client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{TimeoutSeconds: &listTimeout})
	if err != nil {
		return err
	}
	npp.pods = []*PodInfo{}
	for _, pod := range pods.Items {
		var owner string
		if len(pod.OwnerReferences) != 0 {
			owner = pod.OwnerReferences[0].Kind
		}
		npp.addPod(pod.GetName(), pod.GetNamespace(), owner)
	}
	return nil
}


func (npp *NakedPodsProvider) addPod(name, namespace, owner string) {
	log.Info("Found pod ", name, " in namespace ", namespace, " with owner ", owner)
	npp.pods = append(npp.pods, &PodInfo{Name: name, Namespace: namespace, owner: owner})
}

// func (npp *NakedPodsProvider) deletePod(name, namespace string) {
// 	for index, podInfo := range npp.pods {
// 		if podInfo.Name == name && podInfo.Namespace == namespace {
// 			npp.pods = append(npp.pods[:index], npp.pods[index+1:]...)
// 		}
// 	}
// }

func (npp *NakedPodsProvider) GetData() map[string]interface{} {
	nakedPods := []*PodInfo{}
	for _, podInfo := range npp.pods {
		if podInfo.owner == "" {
			nakedPods = append(nakedPods, podInfo)
		}
	}
	return map[string]interface{}{"nakedPods": nakedPods}
}
