package ruleset

import (
	"github.com/devon-caron/metrifuge/k8s/api"
)

// Rule represents a single processing capturegroup
type RuleSet struct {
	APIVersion string         `json:"apiVersion" yaml:"apiVersion"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   api.Metadata   `json:"metadata" yaml:"metadata"`
	Spec       Spec           `json:"spec" yaml:"spec"`
	Status     map[string]any `json:"status,omitempty" yaml:"status,omitempty"`
}

// RuleSpec contains the capturegroup configuration
type Spec struct {
	Rules []*api.Rule `json:"rules" yaml:"rules"`
}

func (rs RuleSet) GetMetadata() api.Metadata {
	return rs.Metadata
}
