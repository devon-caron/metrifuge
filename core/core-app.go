package core

import (
	"os"
	"strconv"

	"github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/k8s/crd"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/devon-caron/metrifuge/receiver"
	"github.com/sirupsen/logrus"
)

var (
	lr              *receiver.LogReceiver
	log             *logrus.Logger
	pipes           []*crd.Pipe
	rules           []*crd.Rule
	metricExporters []*crd.MetricExporter
	logExporters    []*crd.LogExporter
)

func Start() {
	log = logger.Get()
	log.Info("starting api")
	api.StartApi()
	lr = receiver.GetLogReceiver()

	isK8s, err := strconv.ParseBool(os.Getenv("MF_RUNNING_IN_K8S"))
	if err != nil {
		log.Errorf("failed to parse environment variable MF_RUNNING_IN_K8S:%v", err)
		os.Exit(1)
	}

	pipes = initPipes(isK8s)
	rules = initRules(isK8s)
	metricExporters = initMetricExporters(isK8s)
	logExporters = initLogExporters(isK8s)
}
