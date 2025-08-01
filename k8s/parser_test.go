package k8s

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devon-caron/metrifuge/k8s/crd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadTestFile is a helper function to load test YAML files
func loadTestFile(t *testing.T, filename string) []byte {
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read test file: %s", filename)
	return data
}

func TestParseDocuments(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		docCount int
		err      bool
	}{
		{
			name:     "single document",
			filename: "single_doc.yaml",
			docCount: 1,
			err:      false,
		},
		{
			name:     "multiple documents",
			filename: "multiple_docs.yaml",
			docCount: 2,
			err:      false,
		},
		{
			name:     "invalid yaml",
			filename: "invalid.yaml",
			docCount: 0,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadTestFile(t, tt.filename)
			docs, err := parseDocuments(data)

			if tt.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, docs, tt.docCount)
		})
	}
}

func TestParseRules(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		validate func(t *testing.T, rules []*crd.Rule)
		err      bool
	}{
		{
			name:     "valid rule",
			filename: "valid_rule.yaml",
			validate: func(t *testing.T, rules []*crd.Rule) {
				require.Len(t, rules, 1)
				rule := rules[0]
				assert.Equal(t, "mfrule-name", rule.Metadata.Name)
				assert.Equal(t, "sample-name", rule.Spec.Name)
				assert.Equal(t, "%{WORD:grok-word} %{NUMBER:num1} - %{NUMBER:num2}", rule.Spec.Pattern)
				assert.Equal(t, "conditional", rule.Spec.Action)
				assert.NotNil(t, rule.Spec.Conditional)
			},
			err: false,
		},
		// The current parser doesn't validate rule content, so invalid rules will still parse
		// but might fail later during processing
		{
			name:     "empty rule name",
			filename: "invalid_rule.yaml",
			validate: func(t *testing.T, rules []*crd.Rule) {
				require.Len(t, rules, 1)
				assert.Equal(t, "", rules[0].Spec.Name)
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadTestFile(t, tt.filename)
			rules, err := ParseRules(data)

			if tt.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, rules)
			}
		})
	}
}

func TestParsePipes(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		validate func(t *testing.T, pipes []*crd.Pipe)
		err      bool
	}{
		{
			name:     "valid pipe",
			filename: "valid_pipe.yaml",
			validate: func(t *testing.T, pipes []*crd.Pipe) {
				require.Len(t, pipes, 1)
				pipe := pipes[0]
				assert.Equal(t, "mfpipe-name", pipe.Metadata.Name)
				assert.Equal(t, "sample-name", pipe.Spec.Name)
				assert.Equal(t, "default", pipe.Spec.Source.Namespace)
				assert.Equal(t, "app-deployment-2983a99a7be2-8bd", pipe.Spec.Source.Pod)
				assert.Equal(t, "app-container", pipe.Spec.Source.Container)
				assert.Len(t, pipe.Spec.Rules, 1)
				assert.Equal(t, "sample-rule", pipe.Spec.Rules[0].Name)
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadTestFile(t, tt.filename)
			pipes, err := ParsePipes(data)

			if tt.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, pipes)
			}
		})
	}
}

func TestParseLogExporters(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		validate func(t *testing.T, exporters []*crd.LogExporter)
		err      bool
	}{
		{
			name:     "elasticsearch exporter",
			filename: "log_exporter_es.yaml",
			validate: func(t *testing.T, exporters []*crd.LogExporter) {
				require.Len(t, exporters, 1)
				exporter := exporters[0]
				assert.Equal(t, "mflogexporter-name", exporter.Metadata.Name)
				assert.Equal(t, "sample-name", exporter.Spec.Name)
				assert.Equal(t, "elasticsearch", exporter.Spec.Destination.Type)
				require.NotNil(t, exporter.Spec.Destination.Elasticsearch)
				assert.Equal(t, "http://elasticsearch:9200", exporter.Spec.Destination.Elasticsearch.URL)
				assert.Equal(t, "springboot-logs", exporter.Spec.Destination.Elasticsearch.Index)
			},
			err: false,
		},
		{
			name:     "splunk exporter",
			filename: "log_exporter_splunk.yaml",
			validate: func(t *testing.T, exporters []*crd.LogExporter) {
				require.Len(t, exporters, 1)
				exporter := exporters[0]
				assert.Equal(t, "splunk", exporter.Spec.Destination.Type)
				require.NotNil(t, exporter.Spec.Destination.Splunk)
				assert.Equal(t, "https://splunk.company.com:8088/services/collector",
					exporter.Spec.Destination.Splunk.URL)
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadTestFile(t, tt.filename)
			exporters, err := ParseLogExporters(data)

			if tt.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, exporters)
			}
		})
	}
}

func TestParseMetricExporters(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		validate func(t *testing.T, exporters []*crd.MetricExporter)
		err      bool
	}{
		{
			name:     "honeycomb exporter",
			filename: "metric_exporter_honeycomb.yaml",
			validate: func(t *testing.T, exporters []*crd.MetricExporter) {
				require.Len(t, exporters, 1)
				exporter := exporters[0]
				assert.Equal(t, "honeycomb", exporter.Spec.Destination.Type)
				require.NotNil(t, exporter.Spec.Destination.Honeycomb)
				assert.Equal(t, "abc123", exporter.Spec.Destination.Honeycomb.APIKey)
				assert.Equal(t, "springboot-logs", exporter.Spec.Destination.Honeycomb.Dataset)
			},
			err: false,
		},
		{
			name:     "prometheus exporter",
			filename: "metric_exporter_prometheus.yaml",
			validate: func(t *testing.T, exporters []*crd.MetricExporter) {
				require.Len(t, exporters, 1)
				exporter := exporters[0]
				assert.Equal(t, "prometheus", exporter.Spec.Destination.Type)
				require.NotNil(t, exporter.Spec.Destination.Prometheus)
				assert.Equal(t, "http://prometheus:9090",
					exporter.Spec.Destination.Prometheus.Endpoint)
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadTestFile(t, tt.filename)
			exporters, err := ParseMetricExporters(data)

			if tt.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, exporters)
			}
		})
	}
}
