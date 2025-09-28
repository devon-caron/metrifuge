package core

import (
	"fmt"
	"os"

	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	me "github.com/devon-caron/metrifuge/k8s/api/metric_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/devon-caron/metrifuge/k8s"
)

func initClient(config *rest.Config) (*dynamic.DynamicClient, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

func initRuleSets(isK8s bool, k8sClient *dynamic.DynamicClient) ([]*ruleset.RuleSet, error) {
	if !isK8s {
		ruleFilePath := os.Getenv("MF_RULES_FILEPATH")
		data, err := os.ReadFile(ruleFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read myRuleSets file: %v", err)
		}

		myRuleSets, err := k8s.ParseRules(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse myRuleSets file: %v", err)
		}
		return myRuleSets, nil
	}

	myRuleSets, err := k8s.GetK8sResources[ruleset.RuleSet](k8sClient, "RuleSet", "v1alpha1", "rulesets")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myRuleSets, nil
}

func initLogSources(isK8s bool, k8sClient *dynamic.DynamicClient) ([]*ls.LogSource, error) {
	if !isK8s {
		logSourceFilePath := os.Getenv("MF_LOG_SOURCES_FILEPATH")
		data, err := os.ReadFile(logSourceFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read log myLogSources file: %v", err)
		}

		myLogSources, err := k8s.ParseLogSources(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log myLogSources file: %v", err)
		}
		return myLogSources, nil
	}

	myLogSources, err := k8s.GetK8sResources[ls.LogSource](k8sClient, "LogSource", "v1alpha1", "logsources")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myLogSources, nil
}

func initMetricExporters(isK8s bool, k8sClient *dynamic.DynamicClient) ([]*me.MetricExporter, error) {
	if !isK8s {
		metricExporterFilePath := os.Getenv("MF_METRIC_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(metricExporterFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read metric myMetricExporters file: %v", err)
		}

		myMetricExporters, err := k8s.ParseMetricExporters(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse metric myMetricExporters file: %v", err)
		}
		return myMetricExporters, nil
	}

	myMetricExporters, err := k8s.GetK8sResources[me.MetricExporter](k8sClient, "MetricExporter", "v1alpha1", "metricexporters")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myMetricExporters, nil
}

func initLogExporters(isK8s bool, k8sClient *dynamic.DynamicClient) ([]*le.LogExporter, error) {
	if !isK8s {
		logExporterFilePath := os.Getenv("MF_LOG_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(logExporterFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read log myLogExporters file: %v", err)
		}

		myLogExporters, err := k8s.ParseLogExporters(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log myLogExporters file: %v", err)
		}
		return myLogExporters, nil
	}

	myLogExporters, err := k8s.GetK8sResources[le.LogExporter](k8sClient, "LogExporter", "v1alpha1", "logexporters")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myLogExporters, nil
}
