package resource

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

// ParseKubernetesData parses YAML or JSON data
func ParseKubernetesData(data []byte) ([]*KubernetesResource, error) {
	var resources []*KubernetesResource

	// The sigs.k8s.io/yaml library can handle both YAML and JSON
	// It also handles multi-document YAML files (separated by ---)

	// Split multi-document YAML
	docs := splitYAMLDocuments(data)

	for i, doc := range docs {
		if len(doc) == 0 {
			continue
		}

		var resource *KubernetesResource

		// Use sigs.k8s.io/yaml which handles both YAML and JSON
		if err := yaml.Unmarshal(doc, resource); err != nil {
			return nil, fmt.Errorf("failed to parse document %d: %w", i+1, err)
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// splitYAMLDocuments splits multi-document YAML by --- separators
func splitYAMLDocuments(data []byte) [][]byte {
	var documents [][]byte

	// Create a YAML decoder to handle multiple documents
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))

	for {
		var doc interface{}
		err := decoder.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			// If there's an error, return what we have so far
			break
		}

		// Marshal the document back to bytes
		docBytes, err := yaml.Marshal(doc)
		if err != nil {
			continue
		}

		documents = append(documents, docBytes)
	}

	// If no documents were found, return the original data
	if len(documents) == 0 {
		return [][]byte{data}
	}

	return documents
}
