package providers

type Node struct {
	Metadata NodeMetadata `json:"metadata"`
}

type NodeMetadata struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}
