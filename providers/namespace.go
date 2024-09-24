package providers

type Namespace struct {
	Metadata NamespaceMetadata `json:"metadata"`
}

type NamespaceMetadata struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}
