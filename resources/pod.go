package resources

type Pod struct {
	Metadata PodMetadata `json:"metadata"`
	Spec     PodSpec     `json:"spec"`
}

type PodMetadata struct {
	Name            string                      `json:"name"`
	Namespace       string                      `json:"namespace"`
	Labels          map[string]string           `json:"labels"`
	Annotations     map[string]string           `json:"annotations"`
	OwnerReferences []PodMetadataOwnerReference `json:"ownerReferences"`
}

type PodMetadataOwnerReference struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	ApiVersion string `json:"apiVersion"`
}

type PodSpec struct {
	Containers []PodSpecContainers `json:"containers"`
	NodeName   string              `json:"nodeName"`
}

type PodSpecContainers struct {
	Image     string                     `json:"image"`
	Name      string                     `json:"name"`
	Resources PodSpecContainersResources `json:"resources"`
}

type PodSpecContainersResources struct {
	Limits   PodSpecContainersResourcesItem `json:"limits"`
	Requests PodSpecContainersResourcesItem `json:"requests"`
}

type PodSpecContainersResourcesItem struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}
