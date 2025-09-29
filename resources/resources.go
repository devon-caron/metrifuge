package resources

import (
	"sync"

	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	me "github.com/devon-caron/metrifuge/k8s/api/metric_exporter"
	rs "github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"k8s.io/client-go/rest"
)

var (
	instance *Resources
	once     sync.Once
)

type Resources struct {
	mu              sync.RWMutex
	ruleSets        []*rs.RuleSet
	logSources      []*ls.LogSource
	metricExporters []*me.MetricExporter
	logExporters    []*le.LogExporter
	kubeConfig      *rest.Config
	k8sClient       *api.K8sClientWrapper
}

// GetInstance returns the singleton instance of Resources
func GetInstance() *Resources {
	once.Do(func() {
		instance = &Resources{}
	})
	return instance
}

// Getters with read locks
func (r *Resources) GetRuleSets() []*rs.RuleSet {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ruleSets
}

func (r *Resources) GetLogSources() []*ls.LogSource {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.logSources
}

func (r *Resources) GetMetricExporters() []*me.MetricExporter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.metricExporters
}

func (r *Resources) GetLogExporters() []*le.LogExporter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.logExporters
}

func (r *Resources) GetKubeConfig() *rest.Config {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.kubeConfig
}

func (r *Resources) GetK8sClient() *api.K8sClientWrapper {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.k8sClient
}

// Setters with write locks
func (r *Resources) SetRuleSets(ruleSets []*rs.RuleSet) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ruleSets = ruleSets
}

func (r *Resources) SetLogSources(logSources []*ls.LogSource) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logSources = logSources
}

func (r *Resources) SetMetricExporters(metricExporters []*me.MetricExporter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metricExporters = metricExporters
}

func (r *Resources) SetLogExporters(logExporters []*le.LogExporter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logExporters = logExporters
}

func (r *Resources) SetKubeConfig(kubeConfig *rest.Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.kubeConfig = kubeConfig
}

func (r *Resources) SetK8sClient(k8sClient *api.K8sClientWrapper) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.k8sClient = k8sClient
}
