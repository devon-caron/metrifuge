package core

import (
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

		pipes, err := k8s.ParsePipes(data)
		if err != nil {
			log.Errorf("failed to parse pipe file: %v", err)
			return []*crd.Pipe{}
		}
		return pipes
	}

	return []*crd.Pipe{}
}

func initRules(isK8s bool) []*crd.Rule {
	if !isK8s {
		ruleFilePath := os.Getenv("MF_RULES_FILEPATH")
		data, err := os.ReadFile(ruleFilePath)
		if err != nil {
			log.Errorf("failed to read rules file: %v", err)
			return []*crd.Rule{}
		}

		rules, err := k8s.ParseRules(data)
		if err != nil {
			log.Errorf("failed to parse rules file: %v", err)
			return []*crd.Rule{}
		}
		return rules
	}

	return []*crd.Rule{}
}

func initMetricExporters(isK8s bool) []*crd.MetricExporter {
	if !isK8s {
		metricExporterFilePath := os.Getenv("MF_METRIC_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(metricExporterFilePath)
		if err != nil {
			log.Errorf("failed to read metric exporters file: %v", err)
			return []*crd.MetricExporter{}
		}

		exporters, err := k8s.ParseMetricExporters(data)
		if err != nil {
			log.Errorf("failed to parse metric exporters file: %v", err)
			return []*crd.MetricExporter{}
		}
		return exporters
	}

	return []*crd.MetricExporter{}
}

func initLogExporters(isK8s bool) []*crd.LogExporter {
	if !isK8s {
		logExporterFilePath := os.Getenv("MF_LOG_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(logExporterFilePath)
		if err != nil {
			log.Errorf("failed to read log exporters file: %v", err)
			return []*crd.LogExporter{}
		}

		exporters, err := k8s.ParseLogExporters(data)
		if err != nil {
			log.Errorf("failed to parse log exporters file: %v", err)
			return []*crd.LogExporter{}
		}
		return exporters
	}

	return []*crd.LogExporter{}
}
