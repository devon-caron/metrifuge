package core

import (
	"fmt"
	"os"

	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/crd"
)

func initPipes(isK8s bool) []*crd.Pipe {
	if !isK8s {
		pipeFilePath := os.Getenv("MF_PIPES_FILEPATH")
		data, err := os.ReadFile(pipeFilePath)
		if err != nil {
			log.Errorf("failed to read pipe file: %v", err)
			return []*crd.Pipe{}
		}

		pipe, err := k8s.ParsePipe(data)
		if err != nil {
			log.Errorf("failed to parse pipe file: %v", err)
			return []*crd.Pipe{}
		}
		return []*crd.Pipe{pipe}
	}

	return []*crd.Pipe{}
}

func initRules(isK8s bool) []*crd.Rule {
	if !isK8s {
		ruleFilePath := os.Getenv("MF_RULES_FILEPATH")
		fmt.Println(ruleFilePath)
		return []*crd.Rule{}
	}

	return []*crd.Rule{}
}

func initMetricExporters(isK8s bool) []*crd.MetricExporter {
	if !isK8s {
		metricExporterFilePath := os.Getenv("MF_METRIC_EXPORTERS_FILEPATH")
		fmt.Println(metricExporterFilePath)
		return []*crd.MetricExporter{}
	}

	return []*crd.MetricExporter{}
}

func initLogExporters(isK8s bool) []*crd.LogExporter {
	if !isK8s {
		logExporterFilePath := os.Getenv("MF_LOG_EXPORTERS_FILEPATH")
		fmt.Println(logExporterFilePath)
		return []*crd.LogExporter{}
	}

	return []*crd.LogExporter{}
}
