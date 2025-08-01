package pipe

import "github.com/devon-caron/metrifuge/resources"

// Pipe represents a configuration for collecting and processing logs from a specific source
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Pipe struct {
	APIVersion string             `json:"apiVersion" yaml:"apiVersion"`
	Kind       string             `json:"kind" yaml:"kind"`
	Metadata   resources.Metadata `json:"metadata" yaml:"metadata"`
	Spec       PipeSpec           `json:"spec" yaml:"spec"`
}

// PipeSpec contains the pipe configuration
type PipeSpec struct {
	Name     string              `json:"name" yaml:"name"`
	Selector *resources.Selector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Source   *Source             `json:"source" yaml:"source"`
	RuleRefs []RuleRef           `json:"ruleRefs,omitempty" yaml:"ruleRefs,omitempty"`
}

// RuleRef contains data required for referencing rules by Pipes
type RuleRef struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

// Source defines the log source configuration
type Source struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Pod       string `json:"pod" yaml:"pod"`
	Container string `json:"container" yaml:"container"`
}
