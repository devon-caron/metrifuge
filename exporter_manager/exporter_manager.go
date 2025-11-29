package exporter_manager

import (
	"context"
	"fmt"

	"github.com/devon-caron/metrifuge/exporter_manager/log_exporter_client"
	"github.com/devon-caron/metrifuge/exporter_manager/metric_exporter_client"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"github.com/devon-caron/metrifuge/log_handler"
)

type ExporterManager struct {
	exporters map[string]e.Exporter
	mc        *metric_exporter_client.MetricExporterClient
	lc        *log_exporter_client.LogExporterClient
	lh        *log_handler.LogHandler
}

func (em *ExporterManager) Initialize(ctx context.Context, exporters []e.Exporter, logHandler *log_handler.LogHandler) error {
	em.lh = logHandler
	em.exporters = make(map[string]e.Exporter)
	for _, exporter := range exporters {
		em.exporters[exporter.GetMetadata().Name] = exporter
	}

	// Initialize the clients
	em.mc = &metric_exporter_client.MetricExporterClient{}
	if err := em.mc.Initialize(ctx, exporters); err != nil {
		return fmt.Errorf("failed to initialize metric client: %w", err)
	}

	em.lc = &log_exporter_client.LogExporterClient{}
	if err := em.lc.Initialize(ctx, exporters); err != nil {
		return fmt.Errorf("failed to initialize log client: %w", err)
	}

	return nil
}

// TODO 11/28: when ProcessItems is called, add name and namespace of the logsource to each exported metric via context.
func (em *ExporterManager) ProcessItems(ctx context.Context, items []api.ProcessedDataItem) error {
	for _, item := range items {
		// Send metric if present
		// myCtx := context.WithValue
		if item.Metric != nil {
			if err := em.mc.ExportMetric(ctx, item.Metric); err != nil {
				return fmt.Errorf("failed to send metric: %w", err)
			}
		}

		// Send log if present
		if item.ForwardLog != "" {
			if err := em.lc.ExportLog(ctx, item.ForwardLog); err != nil {
				return fmt.Errorf("failed to export log: %w", err)
			}
		}
	}
	return nil
}

func (em *ExporterManager) initializeConnections() {
	for _, exporter := range em.exporters {
		switch exporter.Spec.Destination.Type {
		case "OtelCollector":
			// Initialize OpenTelemetry connection
			// TODO: Implement OTLP connection logic
		case "Honeycomb":
			// Initialize Honeycomb connection
			// TODO: Implement Honeycomb connection logic

		case "Prometheus":
			// Initialize Prometheus connection
			// TODO: Implement Prometheus connection logic
		case "Elasticsearch":
			// Initialize Elasticsearch connection
			// TODO: Implement Elasticsearch connection logic
		case "Splunk":
			// Initialize Splunk connection
			// TODO: Implement Splunk connection logic
		case "Datadog":
			// Initialize Datadog connection
			// TODO: Implement Datadog connection logic
		case "Loki":
			// Initialize Loki connection
			// TODO: Implement Loki connection logic
		}
	}
}
