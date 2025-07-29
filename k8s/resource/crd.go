package resource

// KubernetesResource represents a generic Kubernetes resource
type KubernetesResource struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   map[string]interface{} `json:"metadata" yaml:"metadata"`
	Spec       map[string]interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
}
