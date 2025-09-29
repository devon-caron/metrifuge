package core

import (
	"fmt"
	"os"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	me "github.com/devon-caron/metrifuge/k8s/api/metric_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
)

func updateRuleSets(isK8s bool, k8sClient *api.K8sClientWrapper) ([]*ruleset.RuleSet, error) {
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

	myRuleSets, err := k8s.GetK8sResources[ruleset.RuleSet](k8sClient, global.RULESET_CRD_NAME, "v1alpha1", "rulesets")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myRuleSets, nil
}

func updateLogSources(isK8s bool, k8sClient *api.K8sClientWrapper) ([]*ls.LogSource, error) {
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

	myLogSources, err := k8s.GetK8sResources[ls.LogSource](k8sClient, global.LOGSOURCE_CRD_NAME, "v1alpha1", "logsources")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myLogSources, nil
}

func updateMetricExporters(isK8s bool, k8sClient *api.K8sClientWrapper) ([]*me.MetricExporter, error) {
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

	myMetricExporters, err := k8s.GetK8sResources[me.MetricExporter](k8sClient, global.METRICEXPORTER_CRD_NAME, "v1alpha1", "metricexporters")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myMetricExporters, nil
}

func updateLogExporters(isK8s bool, k8sClient *api.K8sClientWrapper) ([]*le.LogExporter, error) {
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

	myLogExporters, err := k8s.GetK8sResources[le.LogExporter](k8sClient, global.LOGEXPORTER_CRD_NAME, "v1alpha1", "logexporters")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myLogExporters, nil
}
