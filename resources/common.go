package resources

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

// Metadata contains the metadata for a resource definition
type Metadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Selector defines how to select resources
type Selector struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}
