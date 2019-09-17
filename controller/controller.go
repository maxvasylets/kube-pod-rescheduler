package controller

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/maxvasylets/kube-pod-rescheduler/conditions"
	"github.com/maxvasylets/kube-pod-rescheduler/utils"
)

var clientset kubernetes.Interface
var evictionPolicyGroupVersion string

// Run controller
func Run() {
	var err error
	clientset, err = getClient()

	evictionPolicyGroupVersion, err = utils.SupportEviction(clientset)
	if err != nil || len(evictionPolicyGroupVersion) == 0 {
		log.Panic(err.Error())
		return
	}

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("default"))
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
	})

	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func getClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// fallback to kubeconfig
		kubeconfig := filepath.Join("~", ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("The kubeconfig cannot be loaded: %v", err)
		}
	}

	return kubernetes.NewForConfig(config)
}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes pod in the cluster
func onAdd(obj interface{}) {
	// Cast the obj as pod
	pod := obj.(*corev1.Pod)
	fmt.Printf("Add: " + pod.Name + "\n")
	execConditions(pod)
}

func onUpdate(oldObj, newObj interface{}) {
	// Cast the obj as pod
	newPod := newObj.(*corev1.Pod)
	fmt.Printf("Update: " + newPod.Name + "\n")
	execConditions(newPod)
}

func execConditions(pod *corev1.Pod) {
	_, annotationExists := pod.ObjectMeta.Annotations[utils.EvictionAnnotation]
	if !annotationExists {
		fmt.Println("conditions called")
		conditions.RestartsCount(clientset, evictionPolicyGroupVersion, pod)
	}
}
