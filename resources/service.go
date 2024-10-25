package resources

type Service struct {
	Metadata ServiceMetadata `json:"metadata"`
	Spec     ServiceSpec     `json:"spec"`
}

type ServiceMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	// Annotations map[string]string `json:"annotations,omitempty"`
}

type ServiceSpec struct {
	Selector map[string]string `json:"selector,omitempty"`
	Ports    []ServiceSpecPort `json:"ports,omitempty"`
}

type ServiceSpecPort struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}
