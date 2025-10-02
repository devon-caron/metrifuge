package exporter_manager

import (
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	rs "github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"k8s.io/client-go/rest"
)

type ExporterManager struct {
	exporters map[string]e.Exporter
	//otelClient *otel.OpenTelemetryClient
}

func (em *ExporterManager) Initialize(ruleSets []*rs.RuleSet,
	k8sConfig *rest.Config, k8sClient *api.K8sClientWrapper, exporters []e.Exporter) {

	em.exporters = make(map[string]e.Exporter)
	// TODO find a better loop, this looks like shit
	for _, exporter := range exporters {
		for _, ruleSet := range ruleSets {
			if api.MatchLabels(ruleSet.Spec.Selector.MatchLabels, exporter.GetMetadata().Labels) {
				for _, rule := range ruleSet.Spec.Rules {
					exporter.AddRule(rule)
				}
				em.exporters[exporter.GetMetadata().Name] = exporter
			}
		}
	}
}
