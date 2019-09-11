package conditions

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/maxvasylets/kube-pod-resheduler/utils"
)

// Execute condition
func RestartsCount(client clientset.Interface, evictionPolicyGroupVersion string, pod *corev1.Pod) {
	//fmt.Printf("condition: " + pod.Name + "\n")

	for _, c := range pod.Status.ContainerStatuses {
		if c.RestartCount > 1 {
			fmt.Printf("condition: %#v count: %d\n", pod.Name, c.RestartCount)

			success, _ := utils.EvictPod(client, pod, evictionPolicyGroupVersion, false)
			// if err != nil && !errors.IsNotFound(err) {
			// 	log.Panic(err.Error())
			// 	return
			// }
			if success {
				fmt.Printf("Evicted pod: %#v\n", pod.Name)
			} else {
				fmt.Printf("Error when evicting pod: %#v (%#v)\n", pod.Name)
			}

			utils.Annotate(client, pod)
		}
	}
}
