package core

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	exapi "github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/exporter_manager"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
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
	em  *exporter_manager.ExporterManager
)

func Run() {
	logrus.Info("fetching config/env variables...")
	global.InitConfig()
	log = logger.Get()
	log.Info("starting api")
	exapi.StartApi()

	if err := updateResources(); err != nil {
		log.Fatalf("failed to load program resources: %v", err)
		os.Exit(1)
	}

	log.Info("ruleset and exporter resources loaded")
	log.Info("initializing log and inline sources...")

	res := resources.GetInstance()
	lr.Initialize(res.GetLogSources(), log, res.GetKubeConfig(), res.GetK8sClient())
	refresh, err := strconv.Atoi(global.REFRESH_INTERVAL)
	if err != nil {
		log.Warnf("failed to parse environment variable MF_REFRESH_INTERVAL: %v", err)
		log.Warnf("using default value of 60")
		refresh = 60
	}

	log.Info("initializing exporter manager...")

	// First collect all exporters into a single slice
	allExporters := make([]e.Exporter, 0)
	for _, e := range res.GetExporters() {
		allExporters = append(allExporters, *e)
	}

	// Then pass the combined slice
	em.Initialize(res.GetRuleSets(), res.GetKubeConfig(), res.GetK8sClient(), allExporters)

	curRetries := 0
	for {
		if err := updateResources(); err != nil {
			log.Errorf("retrying due to failure to update resources: %v", err)
			time.Sleep(3 * time.Second)
			curRetries++
			if curRetries > 5 {
				log.Fatalf("failed to update resources after 5 retries")
			}
			continue
		}
		curRetries = 0
		time.Sleep(time.Duration(refresh) * time.Second)
		lr.Update(res.GetLogSources(), res.GetK8sClient())
	}
}

func updateResources() error {
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
			return fmt.Errorf("failed to retrieve kubernetes config: %v", err)
		}
		res.SetKubeConfig(kubeConfig)

		k8sClient, err := api.NewK8sClientWrapper(kubeConfig)
		if err != nil {
			return fmt.Errorf("failed to retrieve kubernetes client: %v", err)
		}
		res.SetK8sClient(k8sClient)

		if err = k8s.ValidateResources(kubeConfig); err != nil {
			return fmt.Errorf("failed to validate kubernetes resources: %v", err)
		}
	}

	err = nil

	go func() {
		defer wg.Done()
		if ruleSets, myErr := updateRuleSets(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error updating rulesets 😔: %v}\n", err, myErr)
		} else {
			res.SetRuleSets(ruleSets)
		}
	}()

	go func() {
		defer wg.Done()
		if logSources, myErr := updateLogSources(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error updating log sources 😔: %v}\n", err, myErr)
		} else {
			res.SetLogSources(logSources)
		}
	}()

	go func() {
		defer wg.Done()
		if exporters, myErr := updateExporters(isK8s, res.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error updating exporters 😔: %v}\n", err, myErr)
		} else {
			res.SetExporters(exporters)
		}
	}()

	wg.Wait()

	if err != nil {
		err = fmt.Errorf("failed to initialize resources: \n%v", err)
	}

	return err
}
