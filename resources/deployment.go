package resources

type Deployment struct {
	Metadata DeploymentMetadata `json:"metadata"`
	Spec     DeploymentSpec     `json:"spec"`
}

type DeploymentMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	// Annotations map[string]string `json:"annotations,omitempty"` // for the moment we do not need annotations
}

type DeploymentSpec struct {
	Replicas int32                  `json:"replicas"`
	Selector DeploymentSpecSelector `json:"selector"`
}

type DeploymentSpecSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}
