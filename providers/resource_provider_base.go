package providers

import "k8s.io/apimachinery/pkg/runtime/schema"

type ResourceProviderBase struct {
	Resource      *schema.GroupVersionResource
	ResourcesList *ResourcesList
}
