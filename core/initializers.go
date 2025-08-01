package core

import (
	"os"

	"github.com/devon-caron/metrifuge/k8s"
	le "github.com/devon-caron/metrifuge/resources/log_exporter"
	me "github.com/devon-caron/metrifuge/resources/metric_exporter"
	"github.com/devon-caron/metrifuge/resources/pipe"
	"github.com/devon-caron/metrifuge/resources/rule"
)

func initPipes(isK8s bool) []*pipe.Pipe {
	if !isK8s {
		pipeFilePath := os.Getenv("MF_PIPES_FILEPATH")
		data, err := os.ReadFile(pipeFilePath)
		if err != nil {
			log.Errorf("failed to read pipe file: %v", err)
			return []*pipe.Pipe{}
		}

		pipes, err := k8s.ParsePipes(data)
		if err != nil {
			log.Errorf("failed to parse pipe file: %v", err)
			return []*pipe.Pipe{}
		}
		return pipes
	}

	return []*pipe.Pipe{}
}

func initRules(isK8s bool) []*rule.Rule {
	if !isK8s {
		ruleFilePath := os.Getenv("MF_RULES_FILEPATH")
		data, err := os.ReadFile(ruleFilePath)
		if err != nil {
			log.Errorf("failed to read rules file: %v", err)
			return []*rule.Rule{}
		}

		rules, err := k8s.ParseRules(data)
		if err != nil {
			log.Errorf("failed to parse rules file: %v", err)
			return []*rule.Rule{}
		}
		return rules
	}

	return []*rule.Rule{}
}

func initMetricExporters(isK8s bool) []*me.MetricExporter {
	if !isK8s {
		metricExporterFilePath := os.Getenv("MF_METRIC_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(metricExporterFilePath)
		if err != nil {
			log.Errorf("failed to read metric exporters file: %v", err)
			return []*me.MetricExporter{}
		}

		exporters, err := k8s.ParseMetricExporters(data)
		if err != nil {
			log.Errorf("failed to parse metric exporters file: %v", err)
			return []*me.MetricExporter{}
		}
		return exporters
	}

	return []*me.MetricExporter{}
}

func initLogExporters(isK8s bool) []*le.LogExporter {
	if !isK8s {
		logExporterFilePath := os.Getenv("MF_LOG_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(logExporterFilePath)
		if err != nil {
			log.Errorf("failed to read log exporters file: %v", err)
			return []*le.LogExporter{}
		}

		exporters, err := k8s.ParseLogExporters(data)
		if err != nil {
			log.Errorf("failed to parse log exporters file: %v", err)
			return []*le.LogExporter{}
		}
		return exporters
	}

	return []*le.LogExporter{}
}
