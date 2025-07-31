package k8s

import (
	"fmt"

	"github.com/devon-caron/metrifuge/k8s/crd"
	"gopkg.in/yaml.v3"
)

func ParseRule(data []byte) (*crd.Rule, error) {
	var rule crd.Rule

	// Use sigs.k8s.io/yaml which handles both YAML and JSON
	if err := yaml.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("failed to parse rule document: %w", err)
	}

	return &rule, nil
}

func ParsePipe(data []byte) (*crd.Pipe, error) {
	var pipe crd.Pipe

	// Use sigs.k8s.io/yaml which handles both YAML and JSON
	if err := yaml.Unmarshal(data, &pipe); err != nil {
		return nil, fmt.Errorf("failed to parse pipe document: %w", err)
	}

	return &pipe, nil
}

func ParseLogExporter(data []byte) (*crd.LogExporter, error) {
	var logExporter crd.LogExporter

	// Use sigs.k8s.io/yaml which handles both YAML and JSON
	if err := yaml.Unmarshal(data, &logExporter); err != nil {
		return nil, fmt.Errorf("failed to parse log exporter document: %w", err)
	}

	return &logExporter, nil
}

func ParseMetricExporter(data []byte) (*crd.MetricExporter, error) {
	var metricExporter crd.MetricExporter

	// Use sigs.k8s.io/yaml which handles both YAML and JSON
	if err := yaml.Unmarshal(data, &metricExporter); err != nil {
		return nil, fmt.Errorf("failed to parse metric exporter document: %w", err)
	}

	return &metricExporter, nil
}
