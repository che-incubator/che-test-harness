package deploy

import (
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateEclipseCheCluster create the CR necessary to deploy Eclipse Che
func CreateEclipseCheCluster() *orgv1.CheCluster {
	return &orgv1.CheCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: metadata.Namespace.Name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       CheKind,
			APIVersion: CheAPIVersion,
		},
		Spec: orgv1.CheClusterSpec{
			Server: orgv1.CheClusterSpecServer{
				TlsSupport:     true,
				CheFlavor:      "che",
				CustomCheProperties: map[string]string{
					"CHE_WORKSPACE_SIDECAR_IMAGE__PULL__POLICY": "IfNotPresent",
					"CHE_WORKSPACE_PLUGIN__BROKER_PULL__POLICY": "IfNotPresent",
				},
			},
		},
	}
}
