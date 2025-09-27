package receiver

import (
	"fmt"
	"sync"

	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
)

/**
 * LogReceiver is an internal type that handles the ingestion of logs from
 * all log exporters in the cluster/system.
 */
type LogReceiver struct {
	wg   sync.WaitGroup
	once sync.Once
}

func (lr *LogReceiver) Initialize(initialExporters []*le.LogExporter) {
	lr.once.Do(func() {
		lr.Update(initialExporters)
	})
}

func (lr *LogReceiver) Update(initialExporters []*le.LogExporter) {
	for _, exporter := range initialExporters {
		sourceSpec := exporter.Spec.Source
		if sourceSpec == nil {
			continue
		}

		source := getRawSource(sourceSpec.Type)
		if source == nil {
			continue
		}

		lr.wg.Add(1)
		go func() {
			defer lr.wg.Done()
			lr.receiveLogs(source)
		}()
	}
}

func getRawSource(s string) api.SourceDefinition {
	fmt.Println("unimplemented")
	return &api.PVCSource{}
}

func (lr *LogReceiver) receiveLogs(source api.SourceDefinition) {
	fmt.Println("unimplemented")
}
