package client

import (
	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	groupName = "org.eclipse.che"
)

type K8sClient struct {
	kubeClient *kubernetes.Clientset
}

// NewK8sClient creates kubernetes client wrapper with helper functions and direct access to k8s go client
func NewK8sClient() (*K8sClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	h := &K8sClient{kubeClient: client}
	return h, nil
}

func (c *K8sClient) KubeRest() rest.Interface {
	cfg, _ := config.GetConfig()
	cfg.ContentConfig.GroupVersion = &schema.GroupVersion{Group: groupName, Version: orgv1.SchemeGroupVersion.Version}
	cfg.APIPath = "/apis"
	cfg.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	cfg.UserAgent = rest.DefaultKubernetesUserAgent()
	client, _ := rest.RESTClientFor(cfg)
	return client
}

// Kube returns the clientset for Kubernetes upstream.
func (c *K8sClient) Kube() kubernetes.Interface {
	return c.kubeClient
}
