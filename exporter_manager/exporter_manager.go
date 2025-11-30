package exporter_manager

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/devon-caron/metrifuge/exporter_manager/log_exporter_client"
	"github.com/devon-caron/metrifuge/exporter_manager/metric_exporter_client"
	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	"github.com/sirupsen/logrus"
)

type ExporterManager struct {
	exporters map[string]e.Exporter
	log       *logrus.Logger
	mc        *metric_exporter_client.MetricExporterClient
	lc        *log_exporter_client.LogExporterClient
}

func (em *ExporterManager) Initialize(ctx context.Context, exporters []e.Exporter, log *logrus.Logger) error {
	em.log = log
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
		myCtx := ctx
		if item.LogSourceInfo.Name != "" {
			myCtx = context.WithValue(myCtx, global.SOURCE_NAME_KEY, item.LogSourceInfo.Name)
		}
		if item.LogSourceInfo.Namespace != "" {
			myCtx = context.WithValue(myCtx, global.SOURCE_NAMESPACE_KEY, item.LogSourceInfo.Namespace)
		}
		if item.Metric != nil {
			if err := em.mc.ExportMetric(myCtx, item.Metric); err != nil {
				return fmt.Errorf("failed to send metric: %w", err)
			}
		} else {
			if rand.Intn(1000) == 0 {
				em.log.Debug("empty metric detected (1/1000)")
			}
		}

		// Send log if present
		if item.ForwardLog != "" {
			if err := em.lc.ExportLog(myCtx, item.ForwardLog); err != nil {
				return fmt.Errorf("failed to export log: %w", err)
			}
		} else {
			if rand.Intn(200) == 0 {
				em.log.Debug("blank log detected (1/200)")
			}
		}
	}
	return nil
}
