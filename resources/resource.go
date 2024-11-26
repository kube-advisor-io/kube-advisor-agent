package resources

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

type Resource interface {
	Pod | Deployment | Statefulset | Service | Ingress | Namespace | Node | kyvernov1.ClusterPolicy
}
