package resources

type Resource interface {
	Pod | Deployment | Statefulset | Namespace | Node
}
