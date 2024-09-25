package providers

type Resource interface {
	Pod | Deployment | Statefulset | Namespace | Node
}
