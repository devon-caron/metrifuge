package exporter

import (
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
)

type Exporter struct {
	APIVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind       string       `json:"kind" yaml:"kind"`
	Metadata   api.Metadata `json:"metadata" yaml:"metadata"`
	Spec       ExporterSpec `json:"spec" yaml:"spec"`
	rules      []*api.Rule
}

type ExporterSpec struct {
	Type               string              `json:"type" yaml:"type"`
	MetricExporterSpec MetricExporterSpec  `json:"metricExporter,omitempty" yaml:"metricExporter,omitempty"`
	LogExporterSpec    LogExporterSpec     `json:"logExporter,omitempty" yaml:"logExporter,omitempty"`
	Destination        ExporterDestination `json:"destination" yaml:"destination"`
}

type MetricExporterSpec struct {
	Selector        *api.Selector   `json:"selector,omitempty" yaml:"selector,omitempty"`
	RefreshInterval time.Duration   `json:"refreshInterval,omitempty" yaml:"refreshInterval,omitempty"`
	Source          *api.SourceSpec `json:"source,omitempty" yaml:"source,omitempty"`
}

type LogExporterSpec struct {
	Selector *api.Selector   `json:"selector,omitempty" yaml:"selector,omitempty"`
	Source   *api.SourceSpec `json:"source,omitempty" yaml:"source,omitempty"`
}

type ExporterDestination struct {
	Type          string                   `json:"type" yaml:"type"`
	Honeycomb     *api.HoneycombConfig     `json:"honeycomb,omitempty" yaml:"honeycomb,omitempty"`
	Prometheus    *api.PrometheusConfig    `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Elasticsearch *api.ElasticsearchConfig `json:"elasticsearch,omitempty" yaml:"elasticsearch,omitempty"`
	Splunk        *api.SplunkConfig        `json:"splunk,omitempty" yaml:"splunk,omitempty"`
	Datadog       *api.DatadogConfig       `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Loki          *api.LokiConfig          `json:"loki,omitempty" yaml:"loki,omitempty"`
}

func (e Exporter) GetMetadata() api.Metadata {
	return e.Metadata
}

func (e *Exporter) AddRule(rule *api.Rule) {
	e.rules = append(e.rules, rule)
}

func (e *Exporter) GetDestinationType() string {
	return e.Spec.Destination.Type
}

func (e *Exporter) GetDestinationConfig() any {
	return e.Spec.Destination
}
