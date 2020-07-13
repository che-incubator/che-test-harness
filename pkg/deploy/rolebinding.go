package deploy

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesPullerImageRoleBinding() *rbacv1.RoleBinding{
	return &rbacv1.RoleBinding{
		TypeMeta : metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind: "RoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:"create-daemonset",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: "k8s-image-puller",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     DefaultPullerImageRoleName,
		},
	}
}
