package resources

type Statefulset struct {
	Metadata StatefulsetMetadata `json:"metadata"`
	Spec     StatefulsetSpec     `json:"spec"`
}

type StatefulsetMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	// Annotations map[string]string `json:"annotations,omitempty"` // for the moment we do not need any annotations
}
type StatefulsetSpec struct {
	Replicas int32                  `json:"replicas"`
	Selector DeploymentSpecSelector `json:"selector,omitempty"`
}

type StatefulsetSpecSelector struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}
