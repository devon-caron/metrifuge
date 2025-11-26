package k8s

import (
	"bytes"
	"fmt"

	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/sirupsen/logrus"

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

func ParseRules(data []byte) ([]*ruleset.RuleSet, error) {
	logrus.Info("Parsing rules")
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	rules := make([]*ruleset.RuleSet, 0, len(documents))

	for i, doc := range documents {
		var rule ruleset.RuleSet
		logrus.Debugf("Parsing ruleset document %d, doc: %s", i+1, string(doc))
		if err := yaml.Unmarshal(doc, &rule); err != nil {
			return nil, fmt.Errorf("failed to parse capturegroup document %d: %w", i+1, err)
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func ParseLogSources(data []byte) ([]*ls.LogSource, error) {
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	logSources := make([]*ls.LogSource, 0, len(documents))

	for i, doc := range documents {
		var logSource ls.LogSource
		if err := yaml.Unmarshal(doc, &logSource); err != nil {
			return nil, fmt.Errorf("failed to parse log source document %d: %w", i+1, err)
		}
		logSources = append(logSources, &logSource)
	}

	return logSources, nil
}

func ParseExporters(data []byte) ([]*e.Exporter, error) {
	documents, err := parseDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML documents: %w", err)
	}

	exporters := make([]*e.Exporter, 0, len(documents))

	for i, doc := range documents {
		var exporter e.Exporter
		if err := yaml.Unmarshal(doc, &exporter); err != nil {
			return nil, fmt.Errorf("failed to parse exporter document %d: %w", i+1, err)
		}
		exporters = append(exporters, &exporter)
	}

	return exporters, nil
}
