package api

import "go.opentelemetry.io/otel/attribute"

type MetrifugeK8sResource interface {
	GetMetadata() Metadata
}

// Metadata contains the metadata for a resource definition
type Metadata struct {
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Selector defines how to select resources
type Selector struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}

type ProcessedDataItem struct {
	ForwardLog    string
	Metric        *MetricData
	LogSourceInfo LogSourceInfo
}

type MetricData struct {
	Name       string
	Kind       string
	ValueInt   int64
	ValueFloat float64
	Attributes []attribute.KeyValue
}

func MatchLabels(requiredMatchingLabels map[string]string, labels map[string]string) bool {
	for labelKey, labelValue := range requiredMatchingLabels {
		if labels[labelKey] != labelValue {
			return false
		}
	}
	return true
}
