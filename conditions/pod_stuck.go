package conditions


import (
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// PodStuck
func PodStuck(client clientset.Interface, evictionPolicyGroupVersion string, pod *corev1.Pod) {
  if viper.GetBool("conditions.podStuck.enabled") {
		podStuck(client, evictionPolicyGroupVersion, pod)
	}
}

func podStuck(client clientset.Interface, evictionPolicyGroupVersion string, pod *corev1.Pod) {
   // TODO: implement logic
}
