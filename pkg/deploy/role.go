package deploy

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesPullerImageRole() *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta : metav1.TypeMeta{
			APIVersion:"rbac.authorization.k8s.io/v1",
			Kind: "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: DefaultPullerImageRoleName,
			Labels: map[string]string{
				"app": "kubernetes-image-puller",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"apps"},
				Resources:     []string{"daemonsets"},
				Verbs:         []string{"create", "delete", "watch", "get"},
			},
		},
	}
}
