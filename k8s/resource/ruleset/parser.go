package ruleset

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyInput        = errors.New("empty input")
	ErrInvalidFormat     = errors.New("invalid input format")
	ErrMissingAPIVersion = errors.New("missing required field: apiVersion")
	ErrNoRules           = errors.New("spec must contain at least one rule")
)

// ParseRuleSets parses YAML or JSON data into a slice of RuleSet objects.
// It automatically detects the input format and handles both single and multi-document inputs.
func ParseRuleSets(data []byte) ([]*RuleSet, error) {
	if len(data) == 0 {
		return nil, ErrEmptyInput
	}

	// Trim whitespace to handle cases with just whitespace
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, ErrEmptyInput
	}

	// Try to parse as JSON first (more strict)
	if json.Valid(data) {
		return parseJSON(data)
	}

	// Otherwise, parse as YAML
	rules, err := parseYAML(data)
	if err != nil && errors.Is(err, ErrEmptyInput) {
		// For empty documents, return empty slice instead of error
		return []*RuleSet{}, nil
	}
	return rules, err
}

func parseJSON(data []byte) ([]*RuleSet, error) {
	var rs RuleSet
	if err := yaml.Unmarshal(data, &rs); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := validateRuleSet(&rs); err != nil {
		return nil, err
	}

	return []*RuleSet{&rs}, nil
}

func parseYAML(data []byte) ([]*RuleSet, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var ruleSets []*RuleSet

	for {
		var rs RuleSet
		err := decoder.Decode(&rs)
		if err != nil {
			if err == io.EOF {
				break
			}
			// if strings.Contains(err.Error(), "mapping values are not allowed in this context") {
			// 	return nil, fmt.Errorf("failed to parse YAML; possible empty ruleset: %w", err)
			// }
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}

		// Skip empty documents
		if rs.APIVersion == "" && rs.Metadata == nil && rs.Spec.Rules == nil {
			continue
		}

		if err := validateRuleSet(&rs); err != nil {
			return nil, err
		}

		ruleSets = append(ruleSets, &rs)
	}

	return ruleSets, nil
}

func validateRuleSet(rs *RuleSet) error {
	if rs.APIVersion == "" {
		return ErrMissingAPIVersion
	}

	if rs.Metadata == nil {
		rs.Metadata = make(map[string]interface{})
	}

	if len(rs.Spec.Rules) == 0 {
		return ErrNoRules
	}

	// Validate each rule
	for i, rule := range rs.Spec.Rules {
		if rule.Name == "" {
			return fmt.Errorf("rule at index %d is missing required field: name", i)
		}
		if rule.Spec.Pattern == "" {
			return fmt.Errorf("rule '%s' is missing required field: spec.pattern", rule.Name)
		}
	}

	return nil
}
