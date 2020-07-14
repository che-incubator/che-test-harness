package deploy

const (
	DefaultConfigMapName              = "k8s-image-puller"
	DefaultPullerImageRoleName        = "create-daemonset"
	CheAPIVersion                     = "org.eclipse.che/v1"
	CheKind                           = "CheCluster"
)

func int32Ptr(i int32) *int32 { return &i }
