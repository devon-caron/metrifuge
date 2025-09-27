package logsource

import (
	"fmt"

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
	Type        string         `json:"type" yaml:"type"`
	PVCSource   *api.PVCSource `json:"pvcSource,omitempty" yaml:"pvcSource,omitempty"`
	K8Source    *api.PodSource `json:"k8Source,omitempty" yaml:"k8Source,omitempty"`
	LocalSource *LocalSource   `json:"localSource,omitempty" yaml:"localSource,omitempty"`
	CmdSource   *CmdSource     `json:"cmdSource,omitempty" yaml:"cmdSource,omitempty"`
}

// LocalSource contains the configuration for getting logs from a local file
type LocalSource struct {
	Path string `json:"path" yaml:"path"`
}

// CmdSource contains the configuration for getting logs from a command
// TODO: implement for given pod/container
type CmdSource struct {
	Command string `json:"command" yaml:"command"`
}

func (locs *LocalSource) GetSourceInfo() string {
	return fmt.Sprintf("Local: %s", locs.Path)
}

func (cs *CmdSource) GetSourceInfo() string {
	return fmt.Sprintf("Command: %s", cs.Command)
}
