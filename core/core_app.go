package core

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	exapi "github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/exporter_manager"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	"github.com/devon-caron/metrifuge/resources"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/log_handler"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
	wg  sync.WaitGroup
	lh  *log_handler.LogHandler
	em  *exporter_manager.ExporterManager
	ctx context.Context
)

func Run() {
	logrus.Info("fetching config/env variables...")
	global.InitConfig()
	log = logger.Get()
	ctx = context.Background()
	log.Info("starting api")
	exapi.StartApi()

	if err := validateK8sResources(); err != nil {
		log.Fatalf("failed to validate k8s resource definitions: %v", err)
		// os.Exit(1)
	}

	log.Info("k8s resource definitions validated")
	log.Info("initializing log and inline sources...")

	rsc := resources.GetInstance()
	lh = &log_handler.LogHandler{}
	lh.Initialize(rsc.GetLogSources(), rsc.GetRuleSets(), log, rsc.GetKubeConfig(), rsc.GetK8sClient())
	refresh, err := strconv.Atoi(global.REFRESH_INTERVAL)
	if err != nil {
		log.Warnf("failed to parse environment variable MF_REFRESH_INTERVAL: %v", err)
		log.Warnf("using default value of 60")
		refresh = 60
	}

	log.Info("log/inline sources intialized, starting log handler...")

	go func() {
		curRetries := 0
		for {
			log.Info("updating resources...")
			if err := getResourceUpdates(); err != nil {
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
			lh.Update(rsc.GetLogSources(), rsc.GetK8sClient())
		}
	}()

	log.Info("initializing exporter manager...")

	// Then pass the combined slice
	em = &exporter_manager.ExporterManager{}
	if err := em.Initialize(ctx, rsc.GetExporters(), lh, log); err != nil {
		log.Fatalf("failed to initialize exporter manager: %v", err)
	}

	time.Sleep(30 * time.Minute)
}

func validateK8sResources() error {

	isK8s, err := strconv.ParseBool(global.RUNNING_IN_K8S)
	if err != nil {
		return fmt.Errorf("failed to parse environment variable MF_RUNNING_IN_K8S:%v", err)
	}

	log.Infof("It is %v that the application is running in k8s", isK8s)

	rsc := resources.GetInstance()

	if isK8s {
		kubeConfig, err := k8s.GetKubeConfig()
		if err != nil {
			return fmt.Errorf("failed to retrieve kubernetes config: %v", err)
		}
		rsc.SetKubeConfig(kubeConfig)

		k8sClient, err := api.NewK8sClientWrapper(kubeConfig)
		if err != nil {
			return fmt.Errorf("failed to retrieve kubernetes client: %v", err)
		}
		rsc.SetK8sClient(k8sClient)

		if err = k8s.ValidateResources(kubeConfig); err != nil {
			return fmt.Errorf("failed to validate kubernetes resources: %v", err)
		}
	}

	if err := getResourceUpdates(); err != nil {
		return fmt.Errorf("error updating resources after initialization: %v", err)
	}

	return nil
}

func getResourceUpdates() error {

	log.Info("retrieving resources from cluster...")

	var err error = nil
	wg.Add(3)

	isK8s, err := strconv.ParseBool(global.RUNNING_IN_K8S)
	if err != nil {
		return fmt.Errorf("failed to parse environment variable MF_RUNNING_IN_K8S:%v", err)
	}

	rsc := resources.GetInstance()

	go func() {
		log.Info("retrieving rulesets from cluster...")
		defer wg.Done()
		if ruleSets, myErr := getRuleSetUpdates(isK8s, rsc.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error updating rulesets : %v}\n", err, myErr)
		} else {
			rsc.SetRuleSets(ruleSets)
		}
		log.Info("rulesets updated")
	}()

	go func() {
		log.Info("retrieving log sources from cluster...")
		defer wg.Done()
		if logSources, myErr := getLogSourceUpdates(isK8s, rsc.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error updating log sources : %v}\n", err, myErr)
		} else {
			rsc.SetLogSources(logSources)
		}
		log.Info("log sources updated")
	}()

	go func() {
		log.Info("retrieving exporters from cluster...")
		defer wg.Done()
		if exporters, myErr := getExporterUpdates(isK8s, rsc.GetK8sClient()); myErr != nil {
			err = fmt.Errorf("%v{error updating exporters : %v}\n", err, myErr)
		} else {
			rsc.SetExporters(exporters)
		}
		log.Info("exporters updated")
	}()

	wg.Wait()

	if err != nil {
		err = fmt.Errorf("failed to update resources: \n%v", err)
	}

	log.Info("resources updated")
	return err
}
