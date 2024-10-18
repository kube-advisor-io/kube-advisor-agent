package resources

type Resource interface {
	Pod | Deployment | Statefulset | Service | Ingress | Namespace | Node
}
