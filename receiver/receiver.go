package receiver

import (
	"sync"
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

/**
 * LogReceiver is an internal type that handles the ingestion of logs from
 * all log exporters in the cluster/system.
 */
type LogReceiver struct {
	log             *logrus.Logger
	wg              sync.WaitGroup
	once            sync.Once
	sourceStopChans map[string]chan struct{} // Map of exporter names to their stop channels
	mu              sync.RWMutex             // Protects the exporters map
}

func (lr *LogReceiver) Initialize(initialSources []*ls.LogSource, log *logrus.Logger, kubeConfig *rest.Config, k8sClient *api.K8sClientWrapper) error {
	lr.once.Do(func() {
		lr.log = log
		log.Info("initialized log receiver")
		lr.sourceStopChans = make(map[string]chan struct{})
		log.Info("initialized log receiver sources")
		lr.Update(initialSources, k8sClient)
		log.Info("log receiver updated successfully")
	})

	return nil
}

func (lr *LogReceiver) Update(sources []*ls.LogSource, k8sClient *api.K8sClientWrapper) error {
	// Create a set of current exporter names
	currentSources := make(map[string]bool)
	for _, source := range sources {
		currentSources[source.Metadata.Name] = true
	}

	// Stop and remove any exporters that are no longer present
	lr.mu.Lock()
	for name, stopCh := range lr.sourceStopChans {
		if !currentSources[name] {
			close(stopCh)
			delete(lr.sourceStopChans, name)
		}
	}
	lr.mu.Unlock()

	// Start new exporters or update existing ones
	for _, source := range sources {
		lr.mu.Lock()
		// If exporter already exists, skip or restart it
		if _, exists := lr.sourceStopChans[source.Metadata.Name]; exists {
			lr.mu.Unlock()
			continue
		}

		// Create a new stop channel for this exporter
		stopCh := make(chan struct{})
		lr.sourceStopChans[source.Metadata.Name] = stopCh
		lr.mu.Unlock()

		lr.wg.Add(1)
		go func(name string, src ls.LogSource, ch chan struct{}) {
			defer lr.wg.Done()
			lr.receiveLogs(src, k8sClient, ch)
		}(source.Metadata.Name, *source, stopCh)
	}

	return nil
}

// ShutDown signals all goroutines to stop and waits for them to complete
func (lr *LogReceiver) ShutDown() {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	// Close all stop channels
	for _, stopCh := range lr.sourceStopChans {
		close(stopCh)
	}

	// Clear the exporters map
	lr.sourceStopChans = make(map[string]chan struct{})

	// Wait for all goroutines to complete
	lr.wg.Wait()
}

// StopExporter stops a specific exporter by name
func (lr *LogReceiver) StopExporter(name string) bool {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if stopCh, exists := lr.sourceStopChans[name]; exists {
		close(stopCh)
		delete(lr.sourceStopChans, name)
		return true
	}
	return false
}

func (lr *LogReceiver) receiveLogs(sourceObj ls.LogSource, kClient *api.K8sClientWrapper, stopCh <-chan struct{}) {
	// Create a ticker for periodic processing
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var source api.Source

	switch sourceObj.Spec.Type {
	case "PVCSource":
		source = sourceObj.Spec.Source.PVCSource
	case "PodSource":
		source = sourceObj.Spec.Source.PodSource
	case "LocalSource":
		source = sourceObj.Spec.Source.LocalSource
	case "CmdSource":
		source = sourceObj.Spec.Source.CmdSource
	default:
		lr.log.Errorf("unknown log source type: %s", sourceObj.Spec.Type)
		return
	}

	source.StartLogStream(kClient, nil, stopCh)

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			lr.log.Infof("Processing logs from source: %s\n", source.GetSourceInfo())
			logs := source.GetNewLogs()
			for _, log := range logs {
				lr.log.Debugf("SOURCE: %v", source.GetSourceInfo())
				lr.log.Debugf("LOG: %v", log)
				lr.log.Debugf("LOG PROCESSING STILL NEEDS IMPLEMENTATION")
			}
		}
	}
}
