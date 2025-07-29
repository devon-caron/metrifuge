package ruleset_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devon-caron/metrifuge/k8s/resource/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRuleSet(t *testing.T) {
	// Get the absolute path to the template-rule.yaml file
	yamlPath := filepath.Join("..", "..", "..", "template-rule.yaml")
	yamlData, err := os.ReadFile(yamlPath)
	require.NoError(t, err, "Failed to read template-rule.yaml")

	// Test parsing the YAML
	rules, err := ruleset.ParseRuleSets(yamlData)
	require.NoError(t, err, "Failed to parse RuleSet")
	require.Len(t, rules, 1, "Expected exactly one RuleSet")

	// Verify the parsed RuleSet
	rs := rules[0]
	assert.Equal(t, "metrifuge.com/k8s/v1", rs.APIVersion)

	// Verify metadata
	assert.NotNil(t, rs.Metadata, "Metadata should not be nil")
	assert.Equal(t, "mfruleset-name", rs.Metadata["name"], "Unexpected name in metadata")

	// Verify labels in metadata
	labels, ok := rs.Metadata["labels"].(map[string]interface{})
	require.True(t, ok, "Expected labels to be a map")
	assert.Equal(t, "springboot-logs", labels["app"], "Unexpected app label")

	// Verify spec
	require.NotNil(t, rs.Spec, "Spec should not be nil")
	require.Len(t, rs.Spec.Rules, 1, "Expected exactly one rule")

	// Verify the rule
	rule := rs.Spec.Rules[0]
	assert.Equal(t, "sample-name", rule.Name, "Unexpected rule name")
	assert.Equal(t, "%{WORD:grok-word} %{NUMBER:num1} - %{NUMBER:num2}", rule.Spec.Pattern, "Unexpected pattern")
	assert.Equal(t, "conditional", rule.Spec.Action, "Unexpected action")

	// Verify conditional
	cond := rule.Spec.Conditional
	assert.Equal(t, "num1", cond.Field1.GrokKey, "Unexpected grokKey in field1")
	assert.Equal(t, "Int64", cond.Field1.Type, "Unexpected type in field1")
	assert.Equal(t, "LessThan", cond.Operator, "Unexpected operator")
	assert.Equal(t, "700", cond.Field2.ManualValue, "Unexpected manualValue in field2")
	assert.Equal(t, "Int64", cond.Field2.Type, "Unexpected type in field2")
	assert.Equal(t, "analyze", cond.ActionTrue, "Unexpected actionTrue")
	assert.Equal(t, "discard", cond.ActionFalse, "Unexpected actionFalse")

	// Verify metrics
	require.Len(t, rule.Spec.Metrics, 1, "Expected exactly one metric")
	metric := rule.Spec.Metrics[0]
	assert.Equal(t, "myMetrics.MetricName", metric.Name, "Unexpected metric name")
	assert.Equal(t, "Int64Counter", metric.Kind, "Unexpected metric kind")
	assert.Equal(t, "grok-word", metric.Value.GrokKey, "Unexpected grokKey in metric value")
	assert.Equal(t, "Int64", metric.Value.Type, "Unexpected type in metric value")

	// Test with invalid input
	_, err = ruleset.ParseRuleSets([]byte("invalid yaml"))
	assert.Error(t, err, "Expected error for invalid YAML")
}

func TestParseRuleSet_MultipleDocuments(t *testing.T) {
	// Test with multiple documents in YAML with explicit document separators
	multiDocYAML := `apiVersion: metrifuge.com/k8s/v1
kind: MFRuleSet
metadata:
  name: ruleset-1
  labels:
    app: test-app
spec:
  rules:
    - name: test-rule
      spec:
        pattern: "test"
        action: "discard"
        conditional:
          field1:
            type: "String"
            grokKey: "test"
          operator: "Equals"
          field2:
            type: "String"
            manualValue: "test"
        metrics: []
---
apiVersion: metrifuge.com/k8s/v1
kind: MFRuleSet
metadata:
  name: ruleset-2
  labels:
    app: test-app
spec:
  rules:
    - name: another-rule
      spec:
        pattern: "another-test"
        action: "forward"
        conditional: {}
        metrics: []
`

	rules, err := ruleset.ParseRuleSets([]byte(multiDocYAML))
	require.NoError(t, err, "Failed to parse multi-document YAML")

	if len(rules) != 2 {
		t.Fatalf("Expected 2 RuleSets, got %d", len(rules))
	}

	// Verify first ruleset
	assert.Equal(t, "ruleset-1", rules[0].Metadata["name"])
	assert.Equal(t, 1, len(rules[0].Spec.Rules))
	assert.Equal(t, "test-rule", rules[0].Spec.Rules[0].Name)

	// Verify second ruleset
	assert.Equal(t, "ruleset-2", rules[1].Metadata["name"])
	assert.Equal(t, 1, len(rules[1].Spec.Rules))
	assert.Equal(t, "another-rule", rules[1].Spec.Rules[0].Name)
}

func TestParseRuleSet_JSONInput(t *testing.T) {
	jsonInput := `{
        "apiVersion": "metrifuge.com/k8s/v1",
        "metadata": {"name": "json-ruleset"},
        "spec": {
            "rules": [{
                "name": "json-rule",
                "spec": {
                    "pattern": "test",
                    "action": "discard",
                    "conditional": {
                        "field1": {"type": "String", "grokKey": "test"},
                        "operator": "Equals",
                        "field2": {"type": "String", "manualValue": "test"},
                        "actionTrue": "analyze",
                        "actionFalse": "discard"
                    },
                    "metrics": []
                }
            }]
        }
    }`

	rules, err := ruleset.ParseRuleSets([]byte(jsonInput))
	require.NoError(t, err, "Failed to parse JSON input")
	require.Len(t, rules, 1, "Expected exactly one RuleSet")
	assert.Equal(t, "json-ruleset", rules[0].Metadata["name"])
}

func TestParseRuleSet_EmptyInput(t *testing.T) {
	_, err := ruleset.ParseRuleSets([]byte(""))
	assert.Error(t, err, "Expected error for empty input")
}

func TestParseRuleSet_EmptyDocument(t *testing.T) {
	emptyDoc := "---\n---\n"
	rules, err := ruleset.ParseRuleSets([]byte(emptyDoc))
	require.NoError(t, err, "Should handle empty documents")
	assert.Empty(t, rules, "Should return empty slice for empty documents")
}

func TestParseRuleSet_DeeplyNestedConditional(t *testing.T) {
	deepNested := `apiVersion: metrifuge.com/k8s/v1
metadata:
  name: deep-nested
spec:
  rules:
    - name: nested-rule
      spec:
        pattern: "test"
        action: "conditional"
        conditional:
          field1: {type: "String", grokKey: "level1"}
          operator: "Equals"
          field2: {type: "String", manualValue: "test"}
          actionTrue: "analyze"
          actionFalse: "discard"
          conditional:
            field1: {type: "Int64", grokKey: "level2"}
            operator: "GreaterThan"
            field2: {type: "Int64", manualValue: "100"}
            actionTrue: "forward"
            actionFalse: "discard"
        metrics: []`

	rules, err := ruleset.ParseRuleSets([]byte(deepNested))
	require.NoError(t, err)
	require.Len(t, rules, 1)
	rule := rules[0].Spec.Rules[0]
	assert.NotNil(t, rule.Spec.Conditional.Conditional, "Should handle nested conditionals")
	assert.Equal(t, "level2", rule.Spec.Conditional.Conditional.Field1.GrokKey)
}

func TestParseRuleSet_ComplexMetrics(t *testing.T) {
	complexMetrics := `apiVersion: metrifuge.com/k8s/v1
metadata:
  name: complex-metrics
spec:
  rules:
    - name: metrics-rule
      spec:
        pattern: "test"
        action: "forward"
        conditional: {}
        metrics:
          - name: "test.metric"
            kind: "Int64Counter"
            value: {type: "Int64", grokKey: "count"}
            attributes:
              - name: "status"
                value: {type: "String", grokKey: "status"}
              - name: "endpoint"
                value: {type: "String", manualValue: "/api/test"}`

	rules, err := ruleset.ParseRuleSets([]byte(complexMetrics))
	require.NoError(t, err)
	require.Len(t, rules, 1)
	metrics := rules[0].Spec.Rules[0].Spec.Metrics
	require.Len(t, metrics, 1)
	assert.Equal(t, "test.metric", metrics[0].Name)
	require.Len(t, metrics[0].Attributes, 2)
	assert.Equal(t, "status", metrics[0].Attributes[0].Name)
	assert.Equal(t, "endpoint", metrics[0].Attributes[1].Name)
}
