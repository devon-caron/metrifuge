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

	em.initializeConnections()
}

func (em *ExporterManager) initializeConnections() {
	for _, exporter := range em.exporters {
		switch exporter.Spec.Destination.Type {
		case "otelCollector":
			// Initialize OpenTelemetry connection
			// TODO: Implement OTLP connection logic
		case "honeycomb":
			// Initialize Honeycomb connection
			// TODO: Implement Honeycomb connection logic

		case "prometheus":
			// Initialize Prometheus connection
			// TODO: Implement Prometheus connection logic
		case "elasticsearch":
			// Initialize Elasticsearch connection
			// TODO: Implement Elasticsearch connection logic
		case "splunk":
			// Initialize Splunk connection
			// TODO: Implement Splunk connection logic
		case "datadog":
			// Initialize Datadog connection
			// TODO: Implement Datadog connection logic
		case "loki":
			// Initialize Loki connection
			// TODO: Implement Loki connection logic
		}
	}
}
