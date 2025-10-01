package api

type MetrifugeK8sResource interface {
	GetMetadata() Metadata
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

type SourceSpec struct {
	Type        string       `json:"type" yaml:"type"`
	PVCSource   *PVCSource   `json:"pvcSource,omitempty" yaml:"pvcSource,omitempty"`
	PodSource   *PodSource   `json:"podSource,omitempty" yaml:"podSource,omitempty"`
	LocalSource *LocalSource `json:"localSource,omitempty" yaml:"localSource,omitempty"`
	CmdSource   *CmdSource   `json:"cmdSource,omitempty" yaml:"cmdSource,omitempty"`
}

func MatchLabels(requiredMatchingLabels map[string]string, labels map[string]string) bool {
	for labelKey, labelValue := range requiredMatchingLabels {
		if labels[labelKey] != labelValue {
			return false
		}
	}
	return true
}
