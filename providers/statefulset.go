package providers

type Statefulset struct {
	Metadata StatefulsetMetadata `json:"metadata"` // `mapstructure:"user"`
	Spec     StatefulsetSpec     `json:"spec"`
}

type StatefulsetMetadata struct {
	Name            string                             `json:"name"`
	Namespace       string                             `json:"namespace"`
	Labels          map[string]string                  `json:"labels"`
	Annotations     map[string]string                  `json:"annotations"`
}
type StatefulsetSpec struct {
	Replicas int32                  `json:"replicas"`
	Selector DeploymentSpecSelector `json:"selector"`
}

type StatefulsetSpecSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}
