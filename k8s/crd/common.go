package crd

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
