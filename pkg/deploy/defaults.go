package deploy

const (
	DefaultConfigMapName              = "k8s-image-puller"
	DefaultPullerImageRoleName        = "create-daemonset"
	CodeReadyAPIVersion               = "org.eclipse.che/v1"
	CodeReadyKind                     = "CheCluster"
)

func int32Ptr(i int32) *int32 { return &i }
