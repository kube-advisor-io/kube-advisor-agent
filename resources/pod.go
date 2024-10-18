package resources

type Pod struct {
	Metadata PodMetadata `json:"metadata"`
	Spec     PodSpec     `json:"spec"`
}

type PodMetadata struct {
	Name            string                      `json:"name"`
	Namespace       string                      `json:"namespace"`
	Labels          map[string]string           `json:"labels,omitempty"`
	Annotations     map[string]string           `json:"annotations,omitempty"`
	OwnerReferences []PodMetadataOwnerReference `json:"ownerReferences,omitempty"`
}

type PodMetadataOwnerReference struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	ApiVersion string `json:"apiVersion"`
}

type PodSpec struct {
	Containers []PodSpecContainer `json:"containers"`
	NodeName   string             `json:"nodeName"`
}

type PodSpecContainer struct {
	Image           string                            `json:"image"`
	Name            string                            `json:"name"`
	SecurityContext *PodSpecContainersSecurityContext `json:"securityContext,omitempty"`
	Resources       *PodSpecContainersResources       `json:"resources,omitempty"`
	LivenessProbe   *map[string]interface{}           `json:"livenessProbe,omitempty"`
	ReadinessProbe  *map[string]interface{}           `json:"readinessProbe,omitempty"`
	StartupProbe    *map[string]interface{}           `json:"startupProbe,omitempty"`
}

type PodSpecContainersResources struct {
	Limits   PodSpecContainersResourcesItem `json:"limits,omitempty"`
	Requests PodSpecContainersResourcesItem `json:"requests,omitempty"`
}

type PodSpecContainersResourcesItem struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type PodSpecContainersSecurityContext struct {
	AllowPrivilegeEscalation bool                   `json:"allowPrivilegeEscalation"`
	Capabilites              map[string]interface{} `json:"capabilities,omitempty"`
}
