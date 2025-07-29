package ruleset

type RuleSet struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Metadata   map[string]interface{} `json:"metadata" yaml:"metadata"`
	Spec       struct {
		Rules []Rule `json:"rules" yaml:"rules"`
	} `json:"spec" yaml:"spec"`
}

type Rule struct {
	Name string `json:"name" yaml:"name"`
	Spec struct {
		Pattern     string `json:"pattern" yaml:"pattern"`
		Action      string `json:"action" yaml:"action"`
		Conditional `json:"conditional" yaml:"conditional"`
		Metrics     []Metric `json:"metrics" yaml:"metrics"`
	} `json:"spec" yaml:"spec"`
}

type Conditional struct {
	Field1      Value        `json:"field1" yaml:"field1"`
	Field2      Value        `json:"field2" yaml:"field2"`
	Operator    string       `json:"operator" yaml:"operator"`
	ActionTrue  string       `json:"actionTrue" yaml:"actionTrue"`
	ActionFalse string       `json:"actionFalse" yaml:"actionFalse"`
	Conditional *Conditional `json:"conditional" yaml:"conditional"`
}

type Metric struct {
	Name       string `json:"name" yaml:"name"`
	Kind       string `json:"kind" yaml:"kind"`
	Value      `json:"value" yaml:"value"`
	Attributes []struct {
		Name  string `json:"name" yaml:"name"`
		Value `json:"value" yaml:"value"`
	} `json:"attributes" yaml:"attributes"`
}

type Value struct {
	GrokKey     string `json:"grokKey" yaml:"grokKey"`
	ManualValue string `json:"manualValue" yaml:"manualValue"`
	Type        string `json:"type" yaml:"type"`
}
