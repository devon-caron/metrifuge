package core

import (
	"os"
	"strconv"

	"github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/devon-caron/metrifuge/receiver"
	le "github.com/devon-caron/metrifuge/resources/log_exporter"
	me "github.com/devon-caron/metrifuge/resources/metric_exporter"
	"github.com/devon-caron/metrifuge/resources/pipe"
	"github.com/devon-caron/metrifuge/resources/rule"
	"github.com/sirupsen/logrus"
)

var (
	lr              *receiver.LogReceiver
	log             *logrus.Logger
	pipes           []*pipe.Pipe
	rules           []*rule.Rule
	metricExporters []*me.MetricExporter
	logExporters    []*le.LogExporter
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
