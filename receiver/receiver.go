package receiver

import (
	"fmt"
	"sync"
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

/**
 * LogReceiver is an internal type that handles the ingestion of logs from
 * all log exporters in the cluster/system.
 */
type LogReceiver struct {
	log           *logrus.Logger
	wg            sync.WaitGroup
	once          sync.Once
	exporterChans map[string]chan struct{} // Map of exporter names to their stop channels
	mu            sync.RWMutex             // Protects the exporters map
	KubeConfig    *rest.Config
}

func (lr *LogReceiver) Initialize(initialExporters []*le.LogExporter, log *logrus.Logger, kubeConfig *rest.Config) {
	lr.once.Do(func() {
		lr.log = log
		lr.KubeConfig = kubeConfig
		lr.exporterChans = make(map[string]chan struct{})
		lr.Update(initialExporters)
	})
}

func (lr *LogReceiver) Update(exporters []*le.LogExporter) {
	// Create a set of current exporter names
	currentExporters := make(map[string]bool)
	for _, exporter := range exporters {
		currentExporters[exporter.Metadata.Name] = true
	}

	// Stop and remove any exporters that are no longer present
	lr.mu.Lock()
	for name, stopCh := range lr.exporterChans {
		if !currentExporters[name] {
			close(stopCh)
			delete(lr.exporterChans, name)
		}
	}
	lr.mu.Unlock()

	// Start new exporters or update existing ones
	for _, exporter := range exporters {
		sourceSpec := exporter.Spec.Source
		if sourceSpec == nil {
			continue
		}

		source := getRawSource(sourceSpec)
		if source == nil {
			continue
		}

		lr.mu.Lock()
		// If exporter already exists, skip or restart it
		if _, exists := lr.exporterChans[exporter.Metadata.Name]; exists {
			lr.mu.Unlock()
			continue
		}

		// Create a new stop channel for this exporter
		stopCh := make(chan struct{})
		lr.exporterChans[exporter.Metadata.Name] = stopCh
		lr.mu.Unlock()

		lr.wg.Add(1)
		go func(name string, src api.Source, ch chan struct{}) {
			defer lr.wg.Done()
			lr.receiveLogs(src, ch)
		}(exporter.Metadata.Name, source, stopCh)
	}
}

// ShutDown signals all goroutines to stop and waits for them to complete
func (lr *LogReceiver) ShutDown() {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	// Close all stop channels
	for _, stopCh := range lr.exporterChans {
		close(stopCh)
	}

	// Clear the exporters map
	lr.exporterChans = make(map[string]chan struct{})

	// Wait for all goroutines to complete
	lr.wg.Wait()
}

// StopExporter stops a specific exporter by name
func (lr *LogReceiver) StopExporter(name string) bool {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if stopCh, exists := lr.exporterChans[name]; exists {
		close(stopCh)
		delete(lr.exporterChans, name)
		return true
	}
	return false
}

func getRawSource(s *api.SourceSpec) api.Source {
	panic("unimplemented")
	return &api.PVCSource{}
}

func (lr *LogReceiver) receiveLogs(source api.Source, stopCh <-chan struct{}) {
	// Create a ticker for periodic processing
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	source.StartLogStream(stopCh)

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			fmt.Printf("Processing logs from source: %s\n", source.GetSourceInfo())
			logs := source.GetNewLogs()
			for _, log := range logs {
				fmt.Println(log)
				panic("unimplemented")
			}
		}
	}
}
