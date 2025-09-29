package core

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	exapi "github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	me "github.com/devon-caron/metrifuge/k8s/api/metric_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"k8s.io/client-go/rest"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/devon-caron/metrifuge/receiver"
	"github.com/sirupsen/logrus"
)

var (
	wg              sync.WaitGroup
	lr              *receiver.LogReceiver
	log             *logrus.Logger
	ruleSets        []*ruleset.RuleSet
	logSources      []*ls.LogSource
	metricExporters []*me.MetricExporter
	logExporters    []*le.LogExporter
	KubeConfig      *rest.Config
	K8sClient       *api.K8sClientWrapper
)

func Start() {
	logrus.Info("fetching config/env variables...")
	global.InitConfig()
	log = logger.Get()
	log.Info("starting api")
	exapi.StartApi()

	if err := loadResources(); err != nil {
		log.Fatalf("failed to load program resources: %v", err)
		os.Exit(1)
	}

	log.Info("ruleset and exporter resources loaded")
	log.Info("initializing log and inline sources...")

	lr.Initialize(logSources, log, KubeConfig, K8sClient)
}

func loadResources() error {
	var err error
	wg.Add(4)

	isK8s, err := strconv.ParseBool(global.RUNNING_IN_K8S)
	if err != nil {
		return fmt.Errorf("failed to parse environment variable MF_RUNNING_IN_K8S:%v", err)
	}

	log.Infof("It is %v that the application is running in k8s", isK8s)

	if isK8s {
		if KubeConfig, err = k8s.GetKubeConfig(); err != nil {
			return fmt.Errorf("failed to initialize kubernetes config: %v", err)
		}
		if K8sClient, err = api.NewK8sClientWrapper(KubeConfig); err != nil {
			return fmt.Errorf("failed to initialize kubernetes client: %v", err)
		}
		if err = k8s.ValidateResources(KubeConfig); err != nil {
			return fmt.Errorf("failed to validate kubernetes resources: %v", err)
		}
	}

	go func() {
		var myErr error
		defer wg.Done()
		if ruleSets, myErr = initRuleSets(isK8s, K8sClient); myErr != nil {
			err = fmt.Errorf("%v{error initializing rulesets 😔: %v}\n", err, myErr)
		}
	}()

	go func() {
		var myErr error
		defer wg.Done()
		if logSources, myErr = initLogSources(isK8s, K8sClient); myErr != nil {
			err = fmt.Errorf("%v{error initializing log sources 😔: %v}\n", err, myErr)
		}
	}()

	go func() {
		var myErr error
		defer wg.Done()
		if metricExporters, myErr = initMetricExporters(isK8s, K8sClient); myErr != nil {
			err = fmt.Errorf("%v{error initializing metric exporters 😔: %v}\n", err, myErr)
		}
	}()

	go func() {
		var myErr error
		defer wg.Done()
		if logExporters, myErr = initLogExporters(isK8s, K8sClient); myErr != nil {
			err = fmt.Errorf("%v{error initializing log exporters 😔: %v}\n", err, myErr)
		}
	}()

	wg.Wait()

	if err != nil {
		err = fmt.Errorf("failed to initialize resources: \n%v", err)
	}

	return err
}
