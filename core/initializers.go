package core

import (
	"fmt"
	"os"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
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

func updateExporters(isK8s bool, k8sClient *api.K8sClientWrapper) ([]*e.Exporter, error) {
	if !isK8s {
		exporterFilePath := os.Getenv("MF_EXPORTERS_FILEPATH")
		data, err := os.ReadFile(exporterFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read metric myMetricExporters file: %v", err)
		}

		myExporters, err := k8s.ParseExporters(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse metric myMetricExporters file: %v", err)
		}
		return myExporters, nil
	}

	myExporters, err := k8s.GetK8sResources[e.Exporter](k8sClient, global.EXPORTER_CRD_NAME, "v1alpha1", "exporters")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	return myExporters, nil
}
