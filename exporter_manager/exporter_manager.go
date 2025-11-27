package exporter_manager

import (
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/log_handler"
	"k8s.io/client-go/rest"
)

type ExporterManager struct {
	exporters map[string]e.Exporter
	//otelClient *otel.OpenTelemetryClient
}

func (em *ExporterManager) Initialize(exporters []e.Exporter, logSources []ls.LogSource, k8sConfig *rest.Config, k8sClient *api.K8sClientWrapper, lh *log_handler.LogHandler) {

	em.exporters = make(map[string]e.Exporter)
	// TODO find a better loop, this looks like shit
	for _, exporter := range exporters {
		em.exporters[exporter.GetMetadata().Name] = exporter
	}

}
