package core

import (
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	me "github.com/devon-caron/metrifuge/k8s/api/metric_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/pipe"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"k8s.io/client-go/rest"
	"os"

	"github.com/devon-caron/metrifuge/k8s"
)

func initPipes(isK8s bool) []*pipe.Pipe {
	if !isK8s {
		pipeFilePath := os.Getenv("MF_PIPES_FILEPATH")
		data, err := os.ReadFile(pipeFilePath)
		if err != nil {
			log.Errorf("failed to read pipe file: %v", err)
			return []*pipe.Pipe{}
		}

		myPipes, err := k8s.ParsePipes(data)
		if err != nil {
			log.Errorf("failed to parse pipe file: %v", err)
			return []*pipe.Pipe{}
		}
		return myPipes
	}

	myPipes, err := k8s.GetK8sResources[pipe.Pipe](&rest.Config{}, "Pipe", "v1alpha1", "pipes")
	if err != nil {
		log.Errorf("failed to get k8s resources: %v", err)
		return []*pipe.Pipe{}
	}

	return myPipes
}

func initRuleSets(isK8s bool) []*ruleset.RuleSet {
	if !isK8s {
		ruleFilePath := os.Getenv("MF_RULES_FILEPATH")
		data, err := os.ReadFile(ruleFilePath)
		if err != nil {
			log.Errorf("failed to read myRuleSets file: %v", err)
			return []*ruleset.RuleSet{}
		}

		myRuleSets, err := k8s.ParseRules(data)
		if err != nil {
			log.Errorf("failed to parse myRuleSets file: %v", err)
			return []*ruleset.RuleSet{}
		}
		return myRuleSets
	}

	myRuleSets, err := k8s.GetK8sResources[ruleset.RuleSet](&rest.Config{}, "RuleSet", "v1alpha1", "rulesets")
	if err != nil {
		log.Errorf("failed to get k8s resources: %v", err)
		return []*ruleset.RuleSet{}
	}

	return myRuleSets
}

func initMetricExporters(isK8s bool) []*me.MetricExporter {
	if !isK8s {
		metricExporterFilePath := os.Getenv("MF_METRIC_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(metricExporterFilePath)
		if err != nil {
			log.Errorf("failed to read metric myMetricExporters file: %v", err)
			return []*me.MetricExporter{}
		}

		myMetricExporters, err := k8s.ParseMetricExporters(data)
		if err != nil {
			log.Errorf("failed to parse metric myMetricExporters file: %v", err)
			return []*me.MetricExporter{}
		}
		return myMetricExporters
	}

	myMetricExporters, err := k8s.GetK8sResources[me.MetricExporter](&rest.Config{}, "MetricExporter", "v1alpha1", "metricexporters")
	if err != nil {
		log.Errorf("failed to get k8s resources: %v", err)
		return []*me.MetricExporter{}
	}

	return myMetricExporters
}

func initLogExporters(isK8s bool) []*le.LogExporter {
	if !isK8s {
		logExporterFilePath := os.Getenv("MF_LOG_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(logExporterFilePath)
		if err != nil {
			log.Errorf("failed to read log myLogExporters file: %v", err)
			return []*le.LogExporter{}
		}

		myLogExporters, err := k8s.ParseLogExporters(data)
		if err != nil {
			log.Errorf("failed to parse log myLogExporters file: %v", err)
			return []*le.LogExporter{}
		}
		return myLogExporters
	}

	myLogExporters, err := k8s.GetK8sResources[le.LogExporter](&rest.Config{}, "LogExporter", "v1alpha1", "logexporters")
	if err != nil {
		log.Errorf("failed to get k8s resources: %v", err)
		return []*le.LogExporter{}
	}

	return myLogExporters
}
