package providers

type Deployment struct {
	Metadata DeploymentMetadata `json:"metadata"` // `mapstructure:"user"`
	Spec     DeploymentSpec     `json:"spec"`
}

type DeploymentMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type DeploymentSpec struct {
	Replicas int32                  `json:"replicas"`
	Selector DeploymentSpecSelector `json:"selector"`
}

type DeploymentSpecSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}
