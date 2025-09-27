package api

import "fmt"

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
	MatchLabels map[string]any `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}

type Source struct {
	Type      string `json:"type" yaml:"type"`
	LogSource struct {
		Name string `json:"name" yaml:"name"`
	} `json:"logSource" yaml:"logSource"`
	PVCSource *PVCSource `json:"pvcSource,omitempty" yaml:"pvcSource,omitempty"`
	PodSource *PodSource `json:"podSource,omitempty" yaml:"podSource,omitempty"`
}

type PVCSource struct {
	PVC struct {
		Name string `json:"name" yaml:"name"`
	} `json:"pvc" yaml:"pvc"`
	LogFilePath string `json:"logFilePath" yaml:"logFilePath"`
}

type PodSource struct {
	Pod struct {
		Name      string `json:"name" yaml:"name"`
		Container string `json:"container" yaml:"container"`
	} `json:"pod" yaml:"pod"`
}

type SourceDefinition interface {
	GetSourceInfo() string
}

func (pvc *PVCSource) GetSourceInfo() string {
	return fmt.Sprintf("PVC: %s, Log File Path: %s", pvc.PVC.Name, pvc.LogFilePath)
}

func (pod *PodSource) GetSourceInfo() string {
	return fmt.Sprintf("Pod: %s, Container: %s", pod.Pod.Name, pod.Pod.Container)
}
