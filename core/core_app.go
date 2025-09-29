package core

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	exapi "github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	"github.com/devon-caron/metrifuge/resources"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/devon-caron/metrifuge/receiver"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
	wg  sync.WaitGroup
	lr  *receiver.LogReceiver
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

	res := resources.GetInstance()
	lr.Initialize(res.GetLogSources(), log, res.GetKubeConfig(), res.GetK8sClient())
}

func loadResources() error {
	var err error
	wg.Add(4)

	isK8s, err := strconv.ParseBool(global.RUNNING_IN_K8S)
	if err != nil {
		return fmt.Errorf("failed to parse environment variable MF_RUNNING_IN_K8S:%v", err)
	}

	log.Infof("It is %v that the application is running in k8s", isK8s)

	res := resources.GetInstance()

	if isK8s {
		kubeConfig, err := k8s.GetKubeConfig()
		if err != nil {
			return fmt.Errorf("failed to initialize kubernetes config: %v", err)
		}
		res.SetKubeConfig(kubeConfig)

		k8sClient, err := api.NewK8sClientWrapper(kubeConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize kubernetes client: %v", err)
		}
		res.SetK8sClient(k8sClient)

		if err = k8s.ValidateResources(kubeConfig); err != nil {
			return fmt.Errorf("failed to validate kubernetes resources: %v", err)
		}
	}

	err = nil

	go func() {
		defer wg.Done()
		if ruleSets, myErr := initRuleSets(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error initializing rulesets 😔: %v}\n", err, myErr)
		} else {
			res.SetRuleSets(ruleSets)
		}
	}()

	go func() {
		defer wg.Done()
		if logSources, myErr := initLogSources(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error initializing log sources 😔: %v}\n", err, myErr)
		} else {
			res.SetLogSources(logSources)
		}
	}()

	go func() {
		defer wg.Done()
		if metricExporters, myErr := initMetricExporters(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error initializing metric exporters 😔: %v}\n", err, myErr)
		} else {
			res.SetMetricExporters(metricExporters)
		}
	}()

	go func() {
		defer wg.Done()
		if logExporters, myErr := initLogExporters(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error initializing log exporters 😔: %v}\n", err, myErr)
		} else {
			res.SetLogExporters(logExporters)
		}
	}()

	wg.Wait()

	if err != nil {
		err = fmt.Errorf("failed to initialize resources: \n%v", err)
	}

	return err
}
