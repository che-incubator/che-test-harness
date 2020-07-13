package deploy

import (
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/cmd/operator_osd/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesPullerImageConfigMap() *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: DefaultConfigMapName,
		},
		Data: map[string]string{
			"CACHING_CPU_LIMIT":      ".2",
			"DAEMONSET_NAME":         "kubernetes-image-puller",
			"CACHING_MEMORY_REQUEST": "10Mi",
			"CACHING_INTERVAL_HOURS": "1",
			"CACHING_CPU_REQUEST":    ".05",
			"NAMESPACE":              config.TestHarnessConfig.KubernetesImagePuller.Namespace,
			"NODE_SELECTOR":          "{}",
			"CACHING_MEMORY_LIMIT":   "4000Mi",
			"IMAGES": config.TestHarnessConfig.KubernetesImagePuller.PullerImages,
		},
	}
}
