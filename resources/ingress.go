package resources

type Ingress struct {
	Metadata IngressMetadata `json:"metadata"`
	Spec     IngressSpec     `json:"spec"`
}

type IngressMetadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type IngressSpec struct {
	DefaultBackend *IngressBackend    `json:"defaultBackend,omitempty"`
	Rules          []*IngressSpecRule `json:"rules,omitempty"`
}

type IngressBackend struct {
	Service *IngressBackendService `json:"service,omitempty"`
}

type IngressBackendService struct {
	Name string `json:"name,omitempty"`
}

type IngressSpecRule struct {
	Http *IngressSpecRuleHttp `json:"http,omitempty"`
}

type IngressSpecRuleHttp struct {
	Paths []*IngressSpecRuleHttpPath `json:"paths,omitempty"`
}

type IngressSpecRuleHttpPath struct {
	Backend *IngressBackend `json:"backend,omitempty"`
}
