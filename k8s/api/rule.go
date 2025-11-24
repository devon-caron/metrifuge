package api

type Rule struct {
	Pattern       string           `json:"pattern" yaml:"pattern"`
	Action        string           `json:"action" yaml:"action"` // forward, discard, conditional
	Conditional   *Conditional     `json:"conditional,omitempty" yaml:"conditional,omitempty"`
	CreateMetrics bool             `json:"create_metrics,omitempty" yaml:"create_metrics,omitempty"`
	Metrics       []MetricTemplate `json:"metrics,omitempty" yaml:"metrics,omitempty"`
}

// Conditional defines a condition for capturegroup evaluation
type Conditional struct {
	Field1           FieldValue       `json:"field1" yaml:"field1"`
	Operator         string           `json:"operator" yaml:"operator"` // LessThan, Equals, DoesNotEqual, Exists, DoesNotExist, GreaterThan, GreaterThanOrEqualTo, etc.
	Field2           FieldValue       `json:"field2,omitempty" yaml:"field2,omitempty"`
	ActionTrue       string           `json:"actionTrue" yaml:"actionTrue"`
	ActionFalse      string           `json:"actionFalse" yaml:"actionFalse"`
	MetricsTrue      []MetricTemplate `json:"metricsTrue,omitempty" yaml:"metricsTrue,omitempty"`
	MetricsFalse     []MetricTemplate `json:"metricsFalse,omitempty" yaml:"metricsFalse,omitempty"`
	ConditionalTrue  *Conditional     `json:"conditionalTrue" yaml:"conditionalTrue"`
	ConditionalFalse *Conditional     `json:"conditionalFalse" yaml:"conditionalFalse"`
}

// FieldValue represents a field value that can come from a grok match or be a manual value
type FieldValue struct {
	Type        string `json:"type" yaml:"type"` // Int64, Float64, String
	GrokKey     string `json:"grokKey,omitempty" yaml:"grokKey,omitempty"`
	ManualValue string `json:"manualValue,omitempty" yaml:"manualValue,omitempty"`
}

// MetricTemplate defines a metric to be emitted
type MetricTemplate struct {
	Name       string      `json:"name" yaml:"name"`
	Kind       string      `json:"kind" yaml:"kind"` // Int64Counter, etc.
	Value      MetricValue `json:"value" yaml:"value"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// MetricValue represents the value of a metric
type MetricValue struct {
	Type        string `json:"type" yaml:"type"` // Int64, Float64
	GrokKey     string `json:"grokKey,omitempty" yaml:"grokKey,omitempty"`
	ManualValue string `json:"manualValue,omitempty" yaml:"manualValue,omitempty"`
}

// Attribute represents a key-value pair for metric attributes
type Attribute struct {
	Key   string     `json:"key" yaml:"key"`
	Value FieldValue `json:"value" yaml:"value"`
}
