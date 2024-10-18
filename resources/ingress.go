package resources

type Ingress struct {
	Metadata DeploymentMetadata `json:"metadata"`
	Spec     DeploymentSpec     `json:"spec"`
}

type IngressMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
