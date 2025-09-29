package metric_exporter

import (
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	rs "github.com/devon-caron/metrifuge/k8s/api/ruleset"
)

// MetricExporter represents a configuration for exporting metrics to various destinations
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MetricExporter struct {
	APIVersion string             `json:"apiVersion" yaml:"apiVersion"`
	Kind       string             `json:"kind" yaml:"kind"`
	Metadata   api.Metadata       `json:"metadata" yaml:"metadata"`
	Spec       MetricExporterSpec `json:"spec" yaml:"spec"`
	ruleSets   []*rs.RuleSet
}

// MetricExporterSpec contains the metric exporter configuration
type MetricExporterSpec struct {
	Name            string                    `json:"name" yaml:"name"`
	Selector        *api.Selector             `json:"selector,omitempty" yaml:"selector,omitempty"`
	RefreshInterval time.Duration             `json:"refreshInterval,omitempty" yaml:"refreshInterval,omitempty"`
	Source          *api.SourceSpec           `json:"source,omitempty" yaml:"source,omitempty"`
	Destination     MetricExporterDestination `json:"destination" yaml:"destination"`
}

// MetricExporterDestination defines the destination configuration for metric exporting
type MetricExporterDestination struct {
	Type       string                `json:"type" yaml:"type"` // honeycomb, prometheus, datadog, etc.
	Honeycomb  *api.HoneycombConfig  `json:"honeycomb,omitempty" yaml:"honeycomb,omitempty"`
	Prometheus *api.PrometheusConfig `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Datadog    *api.DatadogConfig    `json:"datadog,omitempty" yaml:"datadog,omitempty"`
}

func (me MetricExporter) GetMetadata() api.Metadata {
	return me.Metadata
}

func (me *MetricExporter) MatchRuleSets(ruleSets []*rs.RuleSet) {
	me.ruleSets = []*rs.RuleSet{}
	for _, ruleSet := range ruleSets {
		if me.Spec.Selector != nil {
			if me.Spec.Selector.MatchLabels != nil {
				allMatched := true
				for labelKey, labelValue := range me.Spec.Selector.MatchLabels {
					if ruleSet.Metadata.Labels[labelKey] != labelValue {
						allMatched = false
						continue
					}
				}
				if allMatched {
					me.ruleSets = append(me.ruleSets, ruleSet)
				}
			}
		}
	}
}
