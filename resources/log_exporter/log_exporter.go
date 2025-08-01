package log_exporter

import "github.com/devon-caron/metrifuge/resources"

// LogExporter represents a configuration for exporting logs to various destinations
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LogExporter struct {
	APIVersion string             `json:"apiVersion" yaml:"apiVersion"`
	Kind       string             `json:"kind" yaml:"kind"`
	Metadata   resources.Metadata `json:"metadata" yaml:"metadata"`
	Spec       LogExporterSpec    `json:"spec" yaml:"spec"`
}

// LogExporterSpec contains the log exporter configuration
type LogExporterSpec struct {
	Name        string                 `json:"name" yaml:"name"`
	Selector    *resources.Selector    `json:"selector,omitempty" yaml:"selector,omitempty"`
	Destination LogExporterDestination `json:"destination" yaml:"destination"`
}

// LogExporterDestination defines the destination configuration for log exporting
type LogExporterDestination struct {
	Type          string                         `json:"type" yaml:"type"` // elasticsearch, splunk, honeycomb, datadog, loki
	Elasticsearch *resources.ElasticsearchConfig `json:"elasticsearch,omitempty" yaml:"elasticsearch,omitempty"`
	Splunk        *resources.SplunkConfig        `json:"splunk,omitempty" yaml:"splunk,omitempty"`
	Honeycomb     *resources.HoneycombConfig     `json:"honeycomb,omitempty" yaml:"honeycomb,omitempty"`
	Datadog       *resources.DatadogConfig       `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Loki          *resources.LokiConfig          `json:"loki,omitempty" yaml:"loki,omitempty"`
}
