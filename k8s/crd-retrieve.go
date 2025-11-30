package k8s

import (
	"context"
	"fmt"
	"slices"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s/api"
	e "github.com/devon-caron/metrifuge/k8s/api/exporter"
	ls "github.com/devon-caron/metrifuge/k8s/api/log_source"
	rs "github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/devon-caron/metrifuge/logger"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	log     = logger.Get()
	crdList *apiextensionsv1.CustomResourceDefinitionList
)

func GetK8sResources(k8sClient *api.K8sClientWrapper, kind, version, kindPlural string) ([]api.MetrifugeK8sResource, error) {
	// For typed client, we would typically use code-generated clients for CRDs
	// Since we don't have those, we'll continue using the dynamic client for now
	// but we'll get it from the rest config in the wrapper

	// First, create a dynamic client using the config from the wrapper
	dynamicClient, err := dynamic.NewForConfig(k8sClient.Config())
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	group := "metrifuge.com"
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: kindPlural,
	}

	log.Debugf("Looking for resources with GVR: %+v", gvr)
	crdResourceList, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %v", err)
	}

	var resources []api.MetrifugeK8sResource
	for _, crdResource := range crdResourceList.Items {
		// Extract spec as map[string]any
		spec, found, err := unstructured.NestedMap(crdResource.Object, "spec")
		if err != nil {
			fmt.Printf("  Error getting spec: %v\n", err)
			continue
		}
		if !found {
			fmt.Printf("  No spec found\n")
			continue
		}

		log.Debugf("Resource spec: %+v", spec)

		resource, err := getResource(crdResource, kind, spec)
		if err != nil {
			log.Warnf("failed to get resource: %v", err)
			continue
		}

		resources = append(resources, resource)
	}
	return resources, nil
}

func getResource(crdResource unstructured.Unstructured, kind string, spec map[string]interface{}) (api.MetrifugeK8sResource, error) {

	log.Infof("Processing resource: %s/%s, kind: %s", crdResource.GetNamespace(), crdResource.GetName(), kind)

	var resource api.MetrifugeK8sResource
	switch kind {
	case global.RULESET_CRD_NAME:

		log.Infof("Processing RuleSet: %s/%s", crdResource.GetNamespace(), crdResource.GetName())
		myRuleSet, err := getRuleSet(crdResource, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get rule set: %v", err)
		}
		resource = myRuleSet
	case global.LOGSOURCE_CRD_NAME:
		log.Infof("Processing LogSource: %s/%s", crdResource.GetNamespace(), crdResource.GetName())
		myLogSource, err := getLogSource(crdResource, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get log source: %v", err)
		}
		resource = myLogSource
	case global.EXPORTER_CRD_NAME:
		log.Infof("Processing Exporter: %s/%s", crdResource.GetNamespace(), crdResource.GetName())
		myExporter, err := getExporter(crdResource, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get exporter: %v", err)
		}
		resource = myExporter
	}

	log.Infof("resource retrieved successfully: %+v", resource)
	return resource, nil
}

func getExporter(crdExporter unstructured.Unstructured, spec map[string]any) (e.Exporter, error) {
	expType, ok := spec["type"].(string)
	if !ok {
		return e.Exporter{}, fmt.Errorf("failed to get exporter type")
	}

	// TODO: If there ends up being a difference, use this
	// switch expType {
	// case "Log":
	// 	// Handle Log exporter type
	// case "Metric":
	// 	// Handle Metric exporter type
	// default:
	// 	return e.Exporter{}, fmt.Errorf("unsupported exporter type: %s", expType)
	// }

	refreshInterval, found := spec["refreshInterval"].(string)
	if !found {
		return e.Exporter{}, fmt.Errorf("failed to get refreshInterval from spec: %v", spec["refreshInterval"])
	}

	destMap, ok := spec["destination"].(map[string]any)
	if !ok {
		return e.Exporter{}, fmt.Errorf("failed to get destination from spec: %v", spec["destination"])
	}

	destination, err := marshalDestination(destMap)
	if err != nil {
		return e.Exporter{}, fmt.Errorf("failed to get destination from spec: %v", err)
	}

	logSourceInfo, ok := spec["logSource"].(map[string]any)
	if !ok {
		return e.Exporter{}, fmt.Errorf("failed to get logSource from spec: %v", spec["logSource"])
	}

	lsName, ok := logSourceInfo["name"].(string)
	if !ok {
		return e.Exporter{}, fmt.Errorf("failed to get logSource name from spec: %v", logSourceInfo["name"])
	}

	lsNamespace, ok := logSourceInfo["namespace"].(string)
	if !ok {
		return e.Exporter{}, fmt.Errorf("failed to get logSource namespace from spec: %v", logSourceInfo["namespace"])
	}

	var myExporter = e.Exporter{
		APIVersion: crdExporter.GetAPIVersion(),
		Kind:       crdExporter.GetKind(),
		Metadata: api.Metadata{
			Name:      crdExporter.GetName(),
			Namespace: crdExporter.GetNamespace(),
		},
		Spec: e.ExporterSpec{
			Type:            expType,
			RefreshInterval: refreshInterval,
			Destination:     destination,
			LogSource: api.LogSourceInfo{
				Name:      lsName,
				Namespace: lsNamespace,
			},
		},
	}

	log.Infof("Exporter spec: %+v", myExporter.Spec)
	log.Infof("Exporter log source: %+v", myExporter.Spec.LogSource)
	log.Infof("Exporter log source name: %s", myExporter.Spec.LogSource.Name)
	log.Infof("Exporter log source namespace: %s", myExporter.Spec.LogSource.Namespace)

	log.Infof("Full exporter object: %+v", myExporter)

	return myExporter, nil
}

func getLogSource(crdLogSource unstructured.Unstructured, spec map[string]any) (ls.LogSource, error) {
	lsSpec, ok := spec["source"].(map[string]any)
	if !ok {
		return ls.LogSource{}, fmt.Errorf("failed to get log source spec: %v", spec)
	}
	lsType := lsSpec["type"].(string)
	switch lsType {
	case "PodSource":
		lsSource, err := marshalPodSource(lsSpec["podSource"].(map[string]any))
		if err != nil {
			return ls.LogSource{}, fmt.Errorf("failed to marshal pod source: %v", err)
		}

		log.Infof("marshaled pod source: %+v", lsSource)
		log.Infof("log source object: %+v", crdLogSource.Object)

		// Extract metadata directly from the Object map
		metadata, found, err := unstructured.NestedMap(crdLogSource.Object, "metadata")
		if !found || err != nil {
			return ls.LogSource{}, fmt.Errorf("failed to get metadata: %v", err)
		}

		name, found, err := unstructured.NestedString(metadata, "name")
		if !found || err != nil {
			return ls.LogSource{}, fmt.Errorf("failed to get name from metadata: %v", err)
		}
		namespace, found, err := unstructured.NestedString(metadata, "namespace")
		if !found || err != nil {
			return ls.LogSource{}, fmt.Errorf("failed to get namespace from metadata: %v", err)
		}
		labels, found, err := unstructured.NestedStringMap(metadata, "labels")
		if !found || err != nil {
			return ls.LogSource{}, fmt.Errorf("failed to get labels from metadata: %v", err)
		}

		log.Infof("Extracted name: '%s', namespace: '%s', labels: %+v", name, namespace, labels)

		return ls.LogSource{
			APIVersion: crdLogSource.GetAPIVersion(),
			Kind:       crdLogSource.GetKind(),
			Metadata: api.Metadata{
				Name:      name,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: ls.LogSourceSpec{
				Type: lsType,
				Source: api.SourceSpec{
					PodSource: lsSource,
				},
			},
		}, nil
	}
	return ls.LogSource{}, fmt.Errorf("unknown log source type: %s", lsType)
}

func getRuleSet(crdRuleSet unstructured.Unstructured, spec map[string]any) (rs.RuleSet, error) {
	rulesList, ok := spec["rules"].([]any)
	if !ok {
		return rs.RuleSet{}, fmt.Errorf("failed to get rules list: %v", spec)
	}

	// Convert []any to []map[string]any
	rulesMaps := make([]map[string]any, 0, len(rulesList))
	for i, r := range rulesList {
		ruleMap, ok := r.(map[string]any)
		if !ok {
			return rs.RuleSet{}, fmt.Errorf("rule at index %d is not a map: %v", i, r)
		}
		rulesMaps = append(rulesMaps, ruleMap)
	}

	myRules, err := getRules(rulesMaps)
	if err != nil {
		return rs.RuleSet{}, fmt.Errorf("failed to get rules: %v", err)
	}

	// Parse selector if present
	var selector *api.Selector
	if selectorMap, ok := spec["selector"].(map[string]any); ok {
		selector, err = marshalSelector(selectorMap)
		if err != nil {
			return rs.RuleSet{}, fmt.Errorf("failed to marshal selector: %v", err)
		}
	}

	// Extract metadata directly from the Object map
	metadata, found, err := unstructured.NestedMap(crdRuleSet.Object, "metadata")
	if !found || err != nil {
		return rs.RuleSet{}, fmt.Errorf("failed to get metadata: %v", err)
	}

	name, _, _ := unstructured.NestedString(metadata, "name")
	namespace, _, _ := unstructured.NestedString(metadata, "namespace")
	labels, _, _ := unstructured.NestedStringMap(metadata, "labels")

	myRuleSet := rs.RuleSet{
		APIVersion: crdRuleSet.GetAPIVersion(),
		Kind:       crdRuleSet.GetKind(),
		Metadata: api.Metadata{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: rs.RuleSetSpec{
			Selector: selector,
			Rules:    myRules,
		},
	}
	return myRuleSet, nil
}

func marshalSelector(selectorMap map[string]any) (*api.Selector, error) {
	matchLabels, ok := selectorMap["matchLabels"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("failed to get matchLabels: %v", selectorMap)
	}

	// Convert map[string]any to map[string]string
	matchLabelsStr := make(map[string]string)
	for k, v := range matchLabels {
		strValue, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("matchLabel value for key %s is not a string: %v", k, v)
		}
		matchLabelsStr[k] = strValue
	}

	return &api.Selector{
		MatchLabels: matchLabelsStr,
	}, nil
}

func marshalPodSource(podSource map[string]any) (*api.PodSource, error) {
	pod, ok := podSource["pod"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("failed to get pod source: %v", podSource)
	}
	myName, ok := pod["name"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get pod name: %v, podSource: %v", pod, podSource)
	}
	myNamespace, ok := pod["namespace"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get pod namespace: %v, podSource: %v", pod, podSource)
	}
	myContainer, ok := pod["container"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get pod container: %v, podSource: %v", pod, podSource)
	}
	var sourceSpec = &api.PodSource{
		Pod: api.Pod{
			Name:      myName,
			Namespace: myNamespace,
			Container: myContainer,
		},
	}

	log.Infof("podSource marshalled successfully: %+v", sourceSpec)

	return sourceSpec, nil
}

func getRules(ruleMaps []map[string]any) ([]*api.Rule, error) {
	var rules []*api.Rule
	for i, ruleMap := range ruleMaps {
		rule, err := getRule(ruleMap)
		if err != nil {
			return nil, fmt.Errorf("failed to get rule at index %d: %v", i, err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func getRule(ruleMap map[string]any) (*api.Rule, error) {

	conditionalMap, ok := ruleMap["conditional"].(map[string]any)
	if !ok {
		conditionalMap = nil
	}

	conditional, err := marshalConditional(conditionalMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal conditional: %v", err)
	}

	pattern, ok := ruleMap["pattern"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get pattern: %v", ruleMap)
	}

	action, ok := ruleMap["action"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get action: %v", ruleMap)
	}

	createMetrics, ok := ruleMap["createMetrics"].(bool)
	if !ok {
		return nil, fmt.Errorf("failed to get createMetrics: %v", ruleMap)
	}

	switch ruleMap["metrics"].(type) {
	case []map[string]any:
		log.Infof("metrics is a slice of maps: %v", ruleMap["metrics"])
	case map[string]any:
		log.Infof("metrics is a map: %v", ruleMap["metrics"])
	default:
		log.Infof("metrics is of type %T", ruleMap["metrics"])
	}

	metricsMap, ok := ruleMap["metrics"].([]any)
	if !ok {
		if ruleMap["metrics"] == nil {
			metricsMap = []any{}
		} else {
			return nil, fmt.Errorf("failed to retrieve metrics: %v", ruleMap)
		}
	}

	metrics, err := marshalMetrics(metricsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metrics: %v", err)
	}

	return &api.Rule{
		Pattern:       pattern,
		Action:        action,
		Conditional:   conditional,
		CreateMetrics: createMetrics,
		Metrics:       metrics,
	}, nil
}

func marshalMetrics(metricsMap []any) ([]api.MetricTemplate, error) {
	var metrics []api.MetricTemplate
	for i, metricMap := range metricsMap {
		metric, err := marshalMetric(metricMap.(map[string]any))
		if err != nil {
			return nil, fmt.Errorf("failed to get metric at index %d: %v", i, err)
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func marshalMetric(metricMap map[string]any) (api.MetricTemplate, error) {
	name, ok := metricMap["name"].(string)
	if !ok {
		return api.MetricTemplate{}, fmt.Errorf("failed to get name: %v", metricMap)
	}
	kind, ok := metricMap["kind"].(string)
	if !ok {
		return api.MetricTemplate{}, fmt.Errorf("failed to get kind: %v", metricMap)
	}
	value, err := marshalMetricValue(metricMap["value"].(map[string]any))
	if err != nil {
		return api.MetricTemplate{}, fmt.Errorf("failed to get value: %v", err)
	}
	var attributes []api.Attribute = nil
	if metricMap["attributes"] != nil {
		attributes, err = marshalAttributes(metricMap["attributes"].([]any))
		if err != nil {
			return api.MetricTemplate{}, fmt.Errorf("failed to get attributes: %v", err)
		}
	}
	return api.MetricTemplate{
		Name:       name,
		Kind:       kind,
		Value:      value,
		Attributes: attributes,
	}, nil
}

func marshalMetricValue(metricValueMap map[string]any) (api.MetricValue, error) {
	metricType, ok := metricValueMap["type"].(string)
	if !ok {
		return api.MetricValue{}, fmt.Errorf("failed to get type: %v", metricValueMap)
	}
	grokKey, ok := metricValueMap["grokKey"].(string)
	if !ok {
		log.Debugf("grok key not found in metric value: %v", metricValueMap)
	}
	manualValue, ok := metricValueMap["manualValue"].(string)
	if !ok {
		log.Debugf("manual value not found in metric value: %v", metricValueMap)
	}
	return api.MetricValue{
		Type:        metricType,
		GrokKey:     grokKey,
		ManualValue: manualValue,
	}, nil
}

func marshalAttributes(attributesMap []any) ([]api.Attribute, error) {
	var attributes []api.Attribute
	for i, attributeMap := range attributesMap {
		attribute, err := marshalAttribute(attributeMap.(map[string]any))
		if err != nil {
			return nil, fmt.Errorf("failed to get attribute at index %d: %v", i, err)
		}
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

func marshalAttribute(attributeMap map[string]any) (api.Attribute, error) {
	key, ok := attributeMap["key"].(string)
	if !ok {
		return api.Attribute{}, fmt.Errorf("failed to get key: %v", attributeMap)
	}
	valueMap, ok := attributeMap["value"].(map[string]any)
	if !ok {
		return api.Attribute{}, fmt.Errorf("failed to get value: %v", attributeMap)
	}
	valueType, ok := valueMap["type"].(string)
	if !ok {
		return api.Attribute{}, fmt.Errorf("failed to get type: %v", valueMap)
	}
	valueGrokKey, ok := valueMap["grokKey"].(string)
	if !ok {
		log.Debugf("grok key not found in attribute value: %v", valueMap)
	}
	valueManualValue, ok := valueMap["manualValue"].(string)
	if !ok {
		log.Debugf("manual value not found in attribute value: %v", valueMap)
	}
	return api.Attribute{
		Key: key,
		Value: api.FieldValue{
			Type:        valueType,
			GrokKey:     valueGrokKey,
			ManualValue: valueManualValue,
		},
	}, nil
}

// marshalConditional marshals a conditional map into an api.Conditional.
// If the conditional map is nil, it returns nil.
func marshalConditional(conditionalMap map[string]any) (*api.Conditional, error) {
	if conditionalMap == nil {
		return nil, nil
	}

	field1, err := marshalFieldValues(conditionalMap["field1"].(map[string]any))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field1: %v", err)
	}
	field2, err := marshalFieldValues(conditionalMap["field2"].(map[string]any))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field2: %v", err)
	}
	operator, ok := conditionalMap["operator"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get operator: %v", conditionalMap)
	}
	actionTrue, ok := conditionalMap["actionTrue"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get actionTrue: %v", conditionalMap)
	}
	actionFalse, ok := conditionalMap["actionFalse"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get actionFalse: %v", conditionalMap)
	}
	var metricsTrue []api.MetricTemplate = nil
	if conditionalMap["metricsTrue"] != nil {
		metricsTrue, err = marshalMetrics(conditionalMap["metricsTrue"].([]any))
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metricsTrue: %v", err)
		}
	}
	var metricsFalse []api.MetricTemplate = nil
	if conditionalMap["metricsFalse"] != nil {
		metricsFalse, err = marshalMetrics(conditionalMap["metricsFalse"].([]any))
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metricsFalse: %v", err)
		}
	}
	var conditionalTrue *api.Conditional = nil
	if conditionalMap["conditionalTrue"] != nil {
		conditionalTrue, err = marshalConditional(conditionalMap["conditionalTrue"].(map[string]any))
		if err != nil {
			return nil, fmt.Errorf("failed to marshal conditionalTrue: %v", err)
		}
	}
	var conditionalFalse *api.Conditional = nil
	if conditionalMap["conditionalFalse"] != nil {
		conditionalFalse, err = marshalConditional(conditionalMap["conditionalFalse"].(map[string]any))
		if err != nil {
			return nil, fmt.Errorf("failed to marshal conditionalFalse: %v", err)
		}
	}
	return &api.Conditional{
		Field1:           field1,
		Operator:         operator,
		Field2:           field2,
		ActionTrue:       actionTrue,
		ActionFalse:      actionFalse,
		MetricsTrue:      metricsTrue,
		MetricsFalse:     metricsFalse,
		ConditionalTrue:  conditionalTrue,
		ConditionalFalse: conditionalFalse,
	}, nil
}

func marshalFieldValues(fieldValueMap map[string]any) (api.FieldValue, error) {
	fvType, ok := fieldValueMap["type"].(string)
	if !ok {
		return api.FieldValue{}, fmt.Errorf("failed to get field value type: %v", fieldValueMap)
	}
	grokKey, ok := fieldValueMap["grokKey"].(string)
	if !ok {
		log.Debugf("field value grok key not present: %v", fieldValueMap)
	}
	manualValue, ok := fieldValueMap["manualValue"].(string)
	if !ok {
		log.Debugf("field value manual value not present: %v", fieldValueMap)
	}
	return api.FieldValue{
		Type:        fvType,
		GrokKey:     grokKey,
		ManualValue: manualValue,
	}, nil
}

func marshalDestination(destMap map[string]any) (api.ExporterDestination, error) {
	destType, ok := destMap["type"].(string)
	if !ok {
		return api.ExporterDestination{}, fmt.Errorf("failed to get destination type: %v", destMap)
	}

	var destination api.ExporterDestination
	switch destType {
	case "OtelCollector":

		otelMap, ok := destMap["otelCollector"].(map[string]any)
		if !ok {
			return api.ExporterDestination{}, fmt.Errorf("failed to get otelCollector config: %v", destMap)
		}

		endpoint, ok := otelMap["endpoint"].(string)
		if !ok {
			return api.ExporterDestination{}, fmt.Errorf("failed to get otelCollector endpoint: %v", otelMap)
		}

		insecure, ok := otelMap["insecure"].(bool)
		if !ok {
			insecure = false // default value
		}

		// Handle OtelCollector destination
		destination = api.ExporterDestination{
			Type: "OtelCollector",
			OtelCollector: &api.OtelCollectorConfig{
				Endpoint: endpoint,
				Insecure: insecure,
			},
		}
	default:
		return api.ExporterDestination{}, fmt.Errorf("unsupported destination type: %s", destType)
	}

	return destination, nil
}

func ValidateResources(restConfig *rest.Config) error {

	var requiredCrdTypes = []string{global.RULESET_CRD_NAME, global.LOGSOURCE_CRD_NAME, global.EXPORTER_CRD_NAME}

	log.Info("creating clientSet...")
	// Create a new clientset which includes the CRD API
	clientset, err := apiextensionsclientset.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %v", err)
	}

	// List all CRDs in the cluster
	crdList, err = clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list CRDs: %v", err)
	}

	var existingCrdTypes []string
	// 3. Print the CRDs
	log.Infof("Found %d Custom Resource Definitions:\n", len(crdList.Items))
	for _, crd := range crdList.Items {
		log.Debugf("  Name: %s", crd.Name)
		log.Debugf("  Group: %s", crd.Spec.Group)
		log.Debugf("  Kind: %s", crd.Spec.Names.Kind)
		log.Debugf("  Version(s): ")
		for i, version := range crd.Spec.Versions {
			if i > 0 {
				fmt.Print(", ")
			}
			log.Debug(version.Name)
		}
		log.Debugf("  Scope: %s", crd.Spec.Scope)
		log.Debug("---")

		if slices.Contains(requiredCrdTypes, crd.Spec.Names.Kind) {
			existingCrdTypes = append(existingCrdTypes, crd.Spec.Names.Kind)
		}
	}

	for _, crdType := range requiredCrdTypes {
		if !slices.Contains(existingCrdTypes, crdType) {
			return fmt.Errorf("required Custom Resource Definition %s not found", crdType)
		}
	}

	log.Info("all required CRDs found, resources validated successfully")
	return nil
}
