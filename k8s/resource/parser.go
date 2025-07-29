package resource

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func ParseRule(data []byte) (*Rule, error) {
	var rule Rule

	// Use sigs.k8s.io/yaml which handles both YAML and JSON
	if err := yaml.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return &rule, nil
}
