package metric_exporter

import (
	"time"

	"github.com/devon-caron/metrifuge/resources"
)

// MetricExporter represents a configuration for exporting metrics to various destinations
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MetricExporter struct {
	APIVersion string             `json:"apiVersion" yaml:"apiVersion"`
	Kind       string             `json:"kind" yaml:"kind"`
	Metadata   resources.Metadata `json:"metadata" yaml:"metadata"`
	Spec       MetricExporterSpec `json:"spec" yaml:"spec"`
}

// MetricExporterSpec contains the metric exporter configuration
type MetricExporterSpec struct {
	Name            string                    `json:"name" yaml:"name"`
	Selector        *resources.Selector       `json:"selector,omitempty" yaml:"selector,omitempty"`
	RefreshInterval time.Duration             `json:"refreshInterval,omitempty" yaml:"refreshInterval,omitempty"`
	Destination     MetricExporterDestination `json:"destination" yaml:"destination"`
}

// MetricExporterDestination defines the destination configuration for metric exporting
type MetricExporterDestination struct {
	Type       string                     `json:"type" yaml:"type"` // honeycomb, prometheus, datadog, etc.
	Honeycomb  *resources.HoneycombConfig `json:"honeycomb,omitempty" yaml:"honeycomb,omitempty"`
	Prometheus *PrometheusConfig          `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Datadog    *resources.DatadogConfig   `json:"datadog,omitempty" yaml:"datadog,omitempty"`
}

// PrometheusConfig contains configuration for Prometheus destination
type PrometheusConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}
