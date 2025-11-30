package api

type ExporterDestination struct {
	Type          string               `json:"type" yaml:"type"`
	Honeycomb     *HoneycombConfig     `json:"honeycomb,omitempty" yaml:"honeycomb,omitempty"`
	Prometheus    *PrometheusConfig    `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Elasticsearch *ElasticsearchConfig `json:"elasticsearch,omitempty" yaml:"elasticsearch,omitempty"`
	Splunk        *SplunkConfig        `json:"splunk,omitempty" yaml:"splunk,omitempty"`
	Datadog       *DatadogConfig       `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Loki          *LokiConfig          `json:"loki,omitempty" yaml:"loki,omitempty"`
	OtelCollector *OtelCollectorConfig `json:"otelCollector,omitempty" yaml:"otelCollector,omitempty"`
}

// HoneycombConfig contains configuration for Honeycomb destination
type HoneycombConfig struct {
	APIKey      string `json:"apiKey" yaml:"apiKey"`
	Dataset     string `json:"dataset" yaml:"dataset"`
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// DatadogConfig contains configuration for Datadog destination
type DatadogConfig struct {
	APIKey  string `json:"apiKey" yaml:"apiKey"`
	Service string `json:"service,omitempty" yaml:"service,omitempty"`
	Source  string `json:"source,omitempty" yaml:"source,omitempty"`
	AppKey  string `json:"appKey,omitempty" yaml:"appKey,omitempty"`
	Site    string `json:"site,omitempty" yaml:"site,omitempty"`
}

// PrometheusConfig contains configuration for Prometheus destination
type PrometheusConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// ElasticsearchConfig contains configuration for Elasticsearch destination
type ElasticsearchConfig struct {
	URL      string `json:"url" yaml:"url"`
	Index    string `json:"index" yaml:"index"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	APIKey   string `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
}

// SplunkConfig contains configuration for Splunk destination
type SplunkConfig struct {
	URL        string `json:"url" yaml:"url"`
	Token      string `json:"token" yaml:"token"`
	Index      string `json:"index,omitempty" yaml:"index,omitempty"`
	Source     string `json:"source,omitempty" yaml:"source,omitempty"`
	SourceType string `json:"sourceType,omitempty" yaml:"sourceType,omitempty"`
}

// LokiConfig contains configuration for Loki destination
type LokiConfig struct {
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

// OtelCollectorConfig contains configuration for OpenTelemetry Collector destination
type OtelCollectorConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Insecure bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
}

type LogSourceInfo struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}
