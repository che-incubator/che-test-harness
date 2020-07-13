package deploy

import (
	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/cmd/operator_osd/config"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/monitors/metadata"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateCodeReadyCluster create the CR necessary to deploy Eclipse Che/ Code Ready Workspaces
func CreateCodeReadyCluster() *orgv1.CheCluster {
	return &orgv1.CheCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.TestHarnessConfig.Flavor,
			Namespace: metadata.Namespace.Name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       CodeReadyKind,
			APIVersion: CodeReadyAPIVersion,
		},
		Spec: orgv1.CheClusterSpec{
			Server: orgv1.CheClusterSpecServer{
				SelfSignedCert: false,
				TlsSupport:     false,
				CheFlavor:      config.TestHarnessConfig.Flavor,
				CustomCheProperties: map[string]string{
					"CHE_WORKSPACE_SIDECAR_IMAGE__PULL__POLICY": "IfNotPresent",
					"CHE_WORKSPACE_PLUGIN__BROKER_PULL__POLICY": "IfNotPresent",
				},
			},
		},
	}
}
