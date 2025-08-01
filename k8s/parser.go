package k8s

import (
	"bytes"
	"fmt"

	"github.com/devon-caron/metrifuge/resources"
	le "github.com/devon-caron/metrifuge/resources/log_exporter"
	"github.com/devon-caron/metrifuge/resources/rule"
	"gopkg.in/yaml.v3"
)

func parseDocuments(data []byte) ([][]byte, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var documents [][]byte

	for {
		var value interface{}
		err := decoder.Decode(&value)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML document: %w", err)
		}

		doc, err := yaml.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal YAML document: %w", err)
		}

		documents = append(documents, doc)
	}

	if len(documents) == 0 {
		return nil, fmt.Errorf("no YAML documents found")
	}

	return documents, nil
}

func ParseRules(data []byte) ([]*rule.Rule, error) {
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	rules := make([]*rule.Rule, 0, len(documents))

	for i, doc := range documents {
		var rule rule.Rule
		if err := yaml.Unmarshal(doc, &rule); err != nil {
			return nil, fmt.Errorf("failed to parse rule document %d: %w", i+1, err)
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func ParsePipes(data []byte) ([]*resources.Pipe, error) {
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	pipes := make([]*resources.Pipe, 0, len(documents))

	for i, doc := range documents {
		var pipe resources.Pipe
		if err := yaml.Unmarshal(doc, &pipe); err != nil {
			return nil, fmt.Errorf("failed to parse pipe document %d: %w", i+1, err)
		}
		pipes = append(pipes, &pipe)
	}

	return pipes, nil
}

func ParseLogExporters(data []byte) ([]*le.LogExporter, error) {
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	exporters := make([]*le.LogExporter, 0, len(documents))

	for i, doc := range documents {
		var exporter le.LogExporter
		if err := yaml.Unmarshal(doc, &exporter); err != nil {
			return nil, fmt.Errorf("failed to parse log exporter document %d: %w", i+1, err)
		}
		exporters = append(exporters, &exporter)
	}

	return exporters, nil
}

func ParseMetricExporters(data []byte) ([]*resources.MetricExporter, error) {
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	exporters := make([]*resources.MetricExporter, 0, len(documents))

	for i, doc := range documents {
		var exporter resources.MetricExporter
		if err := yaml.Unmarshal(doc, &exporter); err != nil {
			return nil, fmt.Errorf("failed to parse metric exporter document %d: %w", i+1, err)
		}
		exporters = append(exporters, &exporter)
	}

	return exporters, nil
}
