package deploy

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesPullerImageConfigMap() *v1.ConfigMap {
	var Images  = "pluginbroker-artifacts=quay.io/eclipse/che-plugin-metadata-broker:v3.2.0;"+
	"che-theia=quay.io/eclipse/che-theia:next;" +
	"java9-sidecar=quay.io/eclipse/che-sidecar-java:8-0cfbacb;" +
	"che-sidecar-dependency-analytics=quay.io/eclipse/che-sidecar-dependency-analytics:0.0.13-a38cb0c;" +
	"che-jwt-proxy=quay.io/eclipse/che-jwtproxy:fd94e60;"+
	"jboss-eap-7=registry.redhat.io/jboss-eap-7/eap73-openjdk8-openshift-rhel7@sha256:f355c9673c09f98c223e73c64ab424dc9f5f756fdeb74a4d33f387411fa27738;"+
	"che-plugin-artifacts-broker=quay.io/eclipse/che-plugin-artifacts-broker:v3.2.0;"+
	"che-theia-endpoint-runtime-binary=quay.io/eclipse/che-theia-endpoint-runtime-binary:next;"+
	"che-sidecar-python=quay.io/eclipse/che-sidecar-python:3.7.3-8f39348;" +
	"che-python=quay.io/eclipse/che-python-3.7:nightly"

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
			"NAMESPACE":              KubernetesImgPullerNS,
			"NODE_SELECTOR":          "{}",
			"CACHING_MEMORY_LIMIT":   "4000Mi",
			"IMAGES": Images,
		},
	}
}
