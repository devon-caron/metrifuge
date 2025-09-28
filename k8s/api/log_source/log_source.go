package logsource

import (
	"github.com/devon-caron/metrifuge/k8s/api"
)

// LogSource represents a configuration for getting logs from a source
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LogSource struct {
	APIVersion string        `json:"apiVersion" yaml:"apiVersion"`
	Kind       string        `json:"kind" yaml:"kind"`
	Metadata   api.Metadata  `json:"metadata" yaml:"metadata"`
	Spec       LogSourceSpec `json:"spec" yaml:"spec"`
}

// LogSourceSpec contains the log source configuration
type LogSourceSpec struct {
	Type   string         `json:"type" yaml:"type"`
	Source api.SourceSpec `json:"source" yaml:"source"`
}

func (ls LogSource) GetMetadata() api.Metadata {
	return ls.Metadata
}
