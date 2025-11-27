package exporter_manager

import (
	"context"
	"fmt"

	"github.com/devon-caron/metrifuge/exporter_manager/log_exporter_client"
	"github.com/devon-caron/metrifuge/exporter_manager/metric_exporter_client"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
)

type ExporterManager struct {
	exporters    map[string]e.Exporter
	metricClient *metric_exporter_client.MetricExporterClient
	logClient    *log_exporter_client.LogExporterClient
}

func (em *ExporterManager) Initialize(ctx context.Context, exporters []e.Exporter) error {
	em.exporters = make(map[string]e.Exporter)
	for _, exporter := range exporters {
		em.exporters[exporter.GetMetadata().Name] = exporter
	}

	// Initialize the clients
	em.metricClient = &metric_exporter_client.MetricExporterClient{}
	if err := em.metricClient.Initialize(ctx, exporters); err != nil {
		return fmt.Errorf("failed to initialize metric client: %w", err)
	}

	em.logClient = &log_exporter_client.LogExporterClient{}
	if err := em.logClient.Initialize(ctx, exporters); err != nil {
		return fmt.Errorf("failed to initialize log client: %w", err)
	}

	return nil
}

func (em *ExporterManager) ProcessItems(ctx context.Context, items []api.ProcessedDataItem) error {
	for _, item := range items {
		// Send metric if present
		if item.Metric != nil {
			if err := em.metricClient.ExportMetric(ctx, item.Metric); err != nil {
				return fmt.Errorf("failed to send metric: %w", err)
			}
		}

		// Send log if present
		if item.ForwardLog != "" {
			if err := em.logClient.ExportLog(ctx, item.ForwardLog); err != nil {
				return fmt.Errorf("failed to export log: %w", err)
			}
		}
	}
	return nil
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
