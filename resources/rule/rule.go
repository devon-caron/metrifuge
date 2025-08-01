package rule

import "github.com/devon-caron/metrifuge/resources"

// Rule represents a single processing rule
type Rule struct {
	APIVersion string             `json:"apiVersion" yaml:"apiVersion"`
	Kind       string             `json:"kind" yaml:"kind"`
	Metadata   resources.Metadata `json:"metadata" yaml:"metadata"`
	Spec       RuleSpec           `json:"spec" yaml:"spec"`
}

// RuleSpec contains the rule configuration
type RuleSpec struct {
	Selector    *resources.Selector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Pattern     string              `json:"pattern" yaml:"pattern"`
	Action      string              `json:"action" yaml:"action"` // forward, discard, analyze, conditional
	Conditional *Conditional        `json:"conditional,omitempty" yaml:"conditional,omitempty"`
	Metrics     []Metric            `json:"metrics,omitempty" yaml:"metrics,omitempty"`
}

// Conditional defines a condition for rule evaluation
type Conditional struct {
	Field1      FieldValue  `json:"field1" yaml:"field1"`
	Operator    string      `json:"operator" yaml:"operator"` // LessThan, Equals, DoesNotEqual, Exists, DoesNotExist, GreaterThan, GreaterThanOrEqualTo, etc.
	Field2      *FieldValue `json:"field2,omitempty" yaml:"field2,omitempty"`
	ActionTrue  string      `json:"actionTrue" yaml:"actionTrue"`
	ActionFalse string      `json:"actionFalse" yaml:"actionFalse"`
}

// FieldValue represents a field value that can come from a grok match or be a manual value
type FieldValue struct {
	Type        string  `json:"type" yaml:"type"` // Int64, Float64, String
	GrokKey     string  `json:"grokKey,omitempty" yaml:"grokKey,omitempty"`
	ManualValue *string `json:"manualValue,omitempty" yaml:"manualValue,omitempty"`
}

// Metric defines a metric to be emitted
type Metric struct {
	Name       string      `json:"name" yaml:"name"`
	Kind       string      `json:"kind" yaml:"kind"` // Int64Counter, etc.
	Value      MetricValue `json:"value" yaml:"value"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// MetricValue represents the value of a metric
type MetricValue struct {
	Type        string  `json:"type" yaml:"type"` // Int64, Float64, String
	GrokKey     string  `json:"grokKey,omitempty" yaml:"grokKey,omitempty"`
	ManualValue *string `json:"manualValue,omitempty" yaml:"manualValue,omitempty"`
}

// Attribute represents a key-value pair for metric attributes
type Attribute struct {
	Name  string      `json:"name" yaml:"name"`
	Value MetricValue `json:"value" yaml:"value"`
}
