package che

import (
	"github.com/che-incubator/che-test-harness/pkg/deploy/context"
	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateEclipseCheCluster create the CR necessary to deploy Eclipse Che
func CreateEclipseCheCluster() *orgv1.CheCluster {
	return &orgv1.CheCluster{
		ObjectMeta: v1.ObjectMeta{
			Name:      context.CheCustomResourceName,
			Namespace: context.TestInstance.CheNamespace,
		},
		TypeMeta: v1.TypeMeta{
			Kind:       context.CheCustomResourceKind,
			APIVersion: context.CheCustomResourceAPIVersion,
		},
		Spec: orgv1.CheClusterSpec{
			Server: orgv1.CheClusterSpecServer{
				CheFlavor: "che",
			},
		},
	}
}
