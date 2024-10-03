package resources

type Node struct {
	Metadata NodeMetadata `json:"metadata"`
}

type NodeMetadata struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
