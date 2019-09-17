package conditions

import (
	"fmt"

	"github.com/maxvasylets/kube-pod-rescheduler/utils"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// RestartsCount is a condition which applies when pod was restarted many times
func RestartsCount(client clientset.Interface, evictionPolicyGroupVersion string, pod *corev1.Pod) {
	if viper.GetBool("conditions.restartsCount.enabled") {
		restartsCount(client, evictionPolicyGroupVersion, pod, viper.GetInt32("conditions.restartsCount.counts"))
	}
}

func restartsCount(client clientset.Interface, evictionPolicyGroupVersion string, pod *corev1.Pod, restartsCount int32) {
	for _, c := range pod.Status.ContainerStatuses {
		if c.RestartCount > restartsCount {
			fmt.Printf("condition: %#v count: %d\n", pod.Name, c.RestartCount)

			success, _ := utils.EvictPod(client, pod, evictionPolicyGroupVersion, false)
			// if err != nil && !errors.IsNotFound(err) {
			// 	log.Panic(err.Error())
			// 	return
			// }
			if success {
				fmt.Printf("Evicted pod: %#v\n", pod.Name)
			} else {
				fmt.Printf("Error when evicting pod: %#v\n", pod.Name)
			}

			utils.Annotate(client, pod)
		}
	}
}
