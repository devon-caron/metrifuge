package log_handler

import (
	"sync"
	"time"

	"github.com/devon-caron/metrifuge/k8s/api"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/devon-caron/metrifuge/log_handler/log_processor"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

/**
 * LogHandler is an internal type that handles the ingestion of logs from
 * all log exporters in the cluster/system.
 */
type LogHandler struct {
	lp              *log_processor.LogProcessor
	log             *logrus.Logger
	wg              sync.WaitGroup
	once            sync.Once
	sourceStopChans map[string]chan struct{} // Map of exporter names to their stop channels
	mu              sync.RWMutex             // Protects the exporters map
}

func (lh *LogHandler) Initialize(initialSources []*ls.LogSource, initialRuleSets []*ruleset.RuleSet, log *logrus.Logger, kubeConfig *rest.Config, k8sClient *api.K8sClientWrapper) error {
	lh.once.Do(func() {
		lh.log = log
		log.Info("initialized log handler")
		lh.sourceStopChans = make(map[string]chan struct{})
		log.Info("initialized log handler sources")
		lh.lp = &log_processor.LogProcessor{}
		lh.lp.Initialize(initialSources, initialRuleSets, log)
		log.Info("log handler updated successfully")
	})

	return nil
}

func (lh *LogHandler) Update(sources []*ls.LogSource, k8sClient *api.K8sClientWrapper) error {
	lh.log.Debug("loghandler update func called")

	// Create a set of current exporter names
	currentSources := make(map[string]bool)
	for _, source := range sources {
		currentSources[source.Metadata.Name] = true
	}

	// Stop and remove any exporters that are no longer present
	lh.mu.Lock()
	for name, stopCh := range lh.sourceStopChans {
		if !currentSources[name] {
			close(stopCh)
			delete(lh.sourceStopChans, name)
		}
	}
	lh.mu.Unlock()

	// Start new exporters or update existing ones
	for _, source := range sources {

		lh.log.Debugf("checking source: %v", source)

		lh.mu.Lock()
		// If exporter already exists, skip or restart it
		if _, exists := lh.sourceStopChans[source.Metadata.Name]; exists {
			lh.mu.Unlock()
			continue
		}

		// Create a new stop channel for this exporter
		stopCh := make(chan struct{})
		lh.sourceStopChans[source.Metadata.Name] = stopCh
		lh.mu.Unlock()

		lh.wg.Add(1)
		go func(name string, src ls.LogSource, ch chan struct{}) {
			defer lh.wg.Done()
			lh.log.Debugf("beginning receipt of logs for source with name %v", name)
			lh.receiveLogs(src, k8sClient, ch)
		}(source.Metadata.Name, *source, stopCh)
	}

	return nil
}

// ShutDown signals all goroutines to stop and waits for them to complete
func (lh *LogHandler) ShutDown() {
	lh.mu.Lock()
	defer lh.mu.Unlock()

	// Close all stop channels
	for _, stopCh := range lh.sourceStopChans {
		close(stopCh)
	}

	// Clear the exporters map
	lh.sourceStopChans = make(map[string]chan struct{})

	// Wait for all goroutines to complete
	lh.wg.Wait()
}

// StopExporter stops a specific exporter by name
func (lh *LogHandler) StopExporter(name string) bool {
	lh.mu.Lock()
	defer lh.mu.Unlock()

	if stopCh, exists := lh.sourceStopChans[name]; exists {
		close(stopCh)
		delete(lh.sourceStopChans, name)
		return true
	}
	return false
}

func (lh *LogHandler) receiveLogs(sourceObj ls.LogSource, kClient *api.K8sClientWrapper, stopCh <-chan struct{}) {
	// Create a ticker for periodic processing
	ticker := time.NewTicker(10 * time.Second)
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
		lh.log.Errorf("unknown log source type: %s", sourceObj.Spec.Type)
		return
	}

	go source.StartLogStream(kClient, nil, stopCh)

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			logs := source.GetNewLogs()
			lh.log.Infof("Processing %v logs from source: %s", len(logs), source.GetSourceInfo())

			logSet, err := lh.lp.FindLogSet(source)
			if err != nil {
				lh.log.Errorf("failed to find log set for source: %v", err)
				continue
			}
			logSet.ProcessLogs(logs)

			// for _, log := range logs {
			// 	lr.log.Infof("received log: %v", log)
			// }
		}
	}
}
