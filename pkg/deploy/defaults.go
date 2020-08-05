package deploy

const (
	DefaultConfigMapName              = "k8s-image-puller"
	DefaultPullerImageRoleName        = "create-daemonset"
	CheAPIVersion                     = "org.eclipse.che/v1"
	CheKind                           = "CheCluster"
	crName                            = "eclipse-che"
	KubernetesImgPullerNS             = "k8s-image-puller"
	K8sIMGPullerContainer             = "quay.io/eclipse/kubernetes-image-puller:latest"
)

func int32Ptr(i int32) *int32 { return &i }
