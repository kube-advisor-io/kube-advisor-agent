package dataproviders

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

var (
	listTimeout      = int64(60)
	lock             = &sync.Mutex{}
	podsListInstance *PodsList
)

type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	owner     string
}

type PodsList struct {
	client                *kubernetes.Clientset
	latestResourceVersion string
	Pods                  []*corev1.Pod
}

func GetPodsListInstance(client *kubernetes.Clientset) *PodsList {
	if podsListInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if podsListInstance == nil {
			podsListInstance = new(PodsList)
			podsListInstance.client = client
			podsListInstance.startWatching()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return podsListInstance
}

func (pl *PodsList) startWatching() {
	err := pl.updatePods()
	if err != nil {
		log.Error(err)
		return
	}
	go pl.watchPods()
}

func (pl *PodsList) watchPods() {
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		log.Info("Starting watching pods...")
		return pl.client.CoreV1().Pods("").Watch(context.Background(), metav1.ListOptions{
			TimeoutSeconds:  &timeOut,
			ResourceVersion: pl.latestResourceVersion,
		})
	}

	watcher, _ := toolsWatch.NewRetryWatcher(pl.latestResourceVersion, &cache.ListWatch{WatchFunc: watchFunc})

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Error:
			fmt.Printf("Error: Object: %v", event.Object)
		case watch.Deleted:
			pod := event.Object.(*corev1.Pod)
			pl.latestResourceVersion = pod.ObjectMeta.ResourceVersion
			pl.deletePod(pod.GetName(), pod.GetNamespace())
		case watch.Added:
			pod := event.Object.(*corev1.Pod)
			pl.latestResourceVersion = pod.ObjectMeta.ResourceVersion
			pl.addPod(pod)
		}
	}
}

func (pl *PodsList) updatePods() error {
	pods, err := pl.client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{TimeoutSeconds: &listTimeout})
	if err != nil {
		return err
	}
	pl.Pods = []*corev1.Pod{}
	for _, pod := range pods.Items {
		pl.addPod(&pod)
	}
	pl.latestResourceVersion = pods.ListMeta.ResourceVersion
	return nil
}

func (pl *PodsList) addPod(pod *corev1.Pod) {
	var owner string
	if len(pod.OwnerReferences) != 0 {
		owner = pod.OwnerReferences[0].Kind
	}
	log.Info("Found pod ", pod.Name, " in namespace ", pod.Namespace, " with owner ", owner)
	pl.Pods = append(pl.Pods, pod)
}

func (pl *PodsList) deletePod(name, namespace string) {
	log.Info("Pod ", name, " in namespace ", namespace, " was deleted")
	for index, pod := range pl.Pods {
		if pod.Name == name && pod.Namespace == namespace {
			pl.Pods = removeFromSlice(pl.Pods, index)
		}
	}
}

func removeFromSlice(s []*corev1.Pod, i int) []*corev1.Pod {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func toPodInfo(pods []*corev1.Pod) []*PodInfo {
	podInfos := []*PodInfo{}
	for _, pod := range pods {
		podInfos = append(podInfos, &PodInfo{Name: pod.Name, Namespace: pod.Namespace})
	}
	return podInfos
}
