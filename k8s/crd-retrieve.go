package k8s

import (
	"context"
	"fmt"
	"slices"

	"github.com/devon-caron/metrifuge/global"
	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
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

func GetK8sResources[Resource api.MetrifugeK8sResource](k8sClient *api.K8sClientWrapper, kind, version, kindPlural string) ([]*Resource, error) {
	// For typed client, we would typically use code-generated clients for CRDs
	// Since we don't have those, we'll continue using the dynamic client for now
	// but we'll get it from the rest config in the wrapper

	// First, create a dynamic client using the config from the wrapper
	dynamicClient, err := dynamic.NewForConfig(k8sClient.Config())
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "metrifuge.com/k8s",
		Version:  version,
		Resource: kindPlural,
	}

	crdResourceList, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var resources []*Resource
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

		resource, err := getResource[Resource](crdResource, kind, spec)
		if err != nil {
			log.Warnf("failed to get resource: %v", err)
			continue
		}

		resources = append(resources, resource)
	}
	return resources, nil
}

func getResource[Resource api.MetrifugeK8sResource](crdResource unstructured.Unstructured, kind string, spec map[string]interface{}) (*Resource, error) {
	var mfK8sCrdNames = []string{global.RULESET_CRD_NAME, global.LOGSOURCE_CRD_NAME, global.LOGEXPORTER_CRD_NAME, global.METRICEXPORTER_CRD_NAME}

	var resource any

	switch kind {
	case mfK8sCrdNames[0]:
		resource = getRuleSet(crdResource, spec)
	case mfK8sCrdNames[1]:
		resource = getLogExporter(crdResource, spec)
	case mfK8sCrdNames[2]:
		resource = getMetricExporter(crdResource, spec)
	}
	r, ok := resource.(Resource)
	if !ok {
		return nil, fmt.Errorf("failed to cast resource: %v", resource)
	}
	return &r, nil
}

func getMetricExporter(crdLogExporter unstructured.Unstructured, spec map[string]any) *le.LogExporter {
	panic("getMetricExporter is not implemented yet")
}

func getLogExporter(crdLogExporter unstructured.Unstructured, spec map[string]any) *le.LogExporter {
	panic("getLogExporter is not implemented yet")
}

func getRuleSet(crdRuleSet unstructured.Unstructured, spec map[string]any) *ruleset.RuleSet {
	rulesMap := spec["rules"].(map[string]any)
	myRules := getRules(rulesMap)

	myRuleSet := &ruleset.RuleSet{
		APIVersion: crdRuleSet.GetAPIVersion(),
		Kind:       crdRuleSet.GetKind(),
		Metadata: api.Metadata{
			Name:      crdRuleSet.GetName(),
			Namespace: crdRuleSet.GetNamespace(),
			Labels:    crdRuleSet.GetLabels(),
		},
		Spec: ruleset.Spec{
			Rules: myRules,
		},
	}
	return myRuleSet
}

func getRules(_ map[string]any) []*ruleset.Rule {
	panic("getRules function not implemented")
}

func ValidateResources(restConfig *rest.Config) error {

	var requiredCrdNames = []string{global.RULESET_CRD_NAME, global.LOGSOURCE_CRD_NAME, global.LOGEXPORTER_CRD_NAME, global.METRICEXPORTER_CRD_NAME}

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

	var existingCrdNames []string
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

		if slices.Contains(requiredCrdNames, crd.Name) {
			existingCrdNames = append(existingCrdNames, crd.Name)
		}
	}

	for _, crdName := range requiredCrdNames {
		if !slices.Contains(existingCrdNames, crdName) {
			return fmt.Errorf("required Custom Resource Definition %s not found", crdName)
		}
	}

	return nil
}
