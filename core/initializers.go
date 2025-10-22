package core

import (
	"fmt"
	"os"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	rs "github.com/devon-caron/metrifuge/k8s/api/ruleset"
)

func updateRuleSets(isK8s bool, k8sClient *api.K8sClientWrapper) ([]*rs.RuleSet, error) {
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

	myResources, err := k8s.GetK8sResources(k8sClient, global.RULESET_CRD_NAME, "v1alpha1", "rulesets")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	myRulesets, err := convertToRulesets(myResources)
	if err != nil {
		return nil, fmt.Errorf("failed to cast resources to rulesets: %v", err)
	}

	return myRulesets, nil
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

	myResources, err := k8s.GetK8sResources(k8sClient, global.LOGSOURCE_CRD_NAME, "v1alpha1", "logsources")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	myLogSources, err := convertToLogSources(myResources)
	if err != nil {
		return nil, fmt.Errorf("failed to cast resources to log sources: %v", err)
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

	myResources, err := k8s.GetK8sResources(k8sClient, global.EXPORTER_CRD_NAME, "v1alpha1", "exporters")
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s resources: %v", err)
	}

	myExporters, err := convertToExporters(myResources)
	if err != nil {
		return nil, fmt.Errorf("failed to cast resources to exporters: %v", err)
	}

	return myExporters, nil
}

func convertToRulesets(resources []api.MetrifugeK8sResource) ([]*rs.RuleSet, error) {
	var myRulesets []*rs.RuleSet
	for i, resource := range resources {
		rs, ok := resource.(*rs.RuleSet)
		if !ok {
			return nil, fmt.Errorf("resource at index %d is not a RuleSet", i)
		}
		myRulesets = append(myRulesets, rs)
	}
	return myRulesets, nil
}

func convertToLogSources(resources []api.MetrifugeK8sResource) ([]*ls.LogSource, error) {
	var myLogSources []*ls.LogSource
	for i, resource := range resources {
		lsrc, ok := resource.(*ls.LogSource)
		if !ok {
			return nil, fmt.Errorf("resource at index %d is not a LogSource", i)
		}
		myLogSources = append(myLogSources, lsrc)
	}
	return myLogSources, nil
}

func convertToExporters(resources []api.MetrifugeK8sResource) ([]*e.Exporter, error) {
	var myExporters []*e.Exporter
	for i, resource := range resources {
		exp, ok := resource.(*e.Exporter)
		if !ok {
			return nil, fmt.Errorf("resource at index %d is not an Exporter", i)
		}
		myExporters = append(myExporters, exp)
	}
	return myExporters, nil
}
