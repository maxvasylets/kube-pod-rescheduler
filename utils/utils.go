package utils

import (
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	EvictionKind        = "Eviction"
	EvictionSubresource = "pods/eviction"
	EvictionAnnotation  = "kube-pod-rescheduler/eviction"
)

// SupportEviction uses Discovery API to find out if the server support eviction subresource
// If support, it will return its groupVersion; Otherwise, it will return ""
func SupportEviction(client clientset.Interface) (string, error) {
	discoveryClient := client.Discovery()
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return "", err
	}
	foundPolicyGroup := false
	var policyGroupVersion string
	for _, group := range groupList.Groups {
		if group.Name == "policy" {
			foundPolicyGroup = true
			policyGroupVersion = group.PreferredVersion.GroupVersion
			break
		}
	}
	if !foundPolicyGroup {
		return "", nil
	}
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion("v1")
	if err != nil {
		return "", err
	}
	for _, resource := range resourceList.APIResources {
		if resource.Name == EvictionSubresource && resource.Kind == EvictionKind {
			return policyGroupVersion, nil
		}
	}
	return "", nil
}

// IsPodExists returns true if the specified pod exists.
// func IsPodExists(client clientset.Interface, podname, namespace string) bool {
// 	_, err := client.CoreV1().Pods(namespace).Get(podname, metav1.GetOptions{})
// 	return !errors.IsNotFound(err)
// }

func Annotate(client clientset.Interface, pod *corev1.Pod) {
	newPod := pod.DeepCopy()
	if newPod.ObjectMeta.Annotations == nil {
		newPod.ObjectMeta.Annotations = make(map[string]string)
	}
	newPod.ObjectMeta.Annotations[EvictionAnnotation] = "true"

	client.CoreV1().Pods(newPod.GetNamespace()).Update(newPod)
}
