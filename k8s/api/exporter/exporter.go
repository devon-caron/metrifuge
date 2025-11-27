package exporter

import (
	"github.com/devon-caron/metrifuge/k8s/api"
)

type Exporter struct {
	APIVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind       string       `json:"kind" yaml:"kind"`
	Metadata   api.Metadata `json:"metadata" yaml:"metadata"`
	Spec       ExporterSpec `json:"spec" yaml:"spec"`
}

type ExporterSpec struct {
	Type            string                  `json:"type" yaml:"type"`
	RefreshInterval string                  `json:"refreshInterval" yaml:"refreshInterval"`
	Destination     api.ExporterDestination `json:"destination" yaml:"destination"`
	LogSource       *api.LogSourceInfo      `json:"logSource" yaml:"logSource"`
}

func (e Exporter) GetMetadata() api.Metadata {
	return e.Metadata
}

func (e *Exporter) GetDestinationType() string {
	return e.Spec.Destination.Type
}

func (e *Exporter) GetDestinationConfig() any {
	return e.Spec.Destination
}
