package k8s

// for listing CRD, go provides client which is different from
// "kubernetes.clientset"
// This clientset will be used to list down the existing CRD
import (
	"context"
	"fmt"
	"slices"

	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/devon-caron/metrifuge/logger"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	log     = logger.Get()
	crdList *v1.CustomResourceDefinitionList
)

func GetK8sResources[Resource api.MetrifugeK8sResource](k8sClient *dynamic.DynamicClient, kind, version, kindPlural string) ([]*Resource, error) {

	gvr := schema.GroupVersionResource{
		Group:    "metrifuge.com/k8s",
		Version:  version,
		Resource: kindPlural,
	}

	crdResourceList, err := k8sClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
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
	var mfK8sCrdNames = []string{"RuleSet", "LogExporter", "MetricExporter"}

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

func validateResources(restConfig *rest.Config) error {

	var requiredCrdNames = []string{"RuleSet", "LogExporter", "MetricExporter"}

	log.Info("creating clientSet...")
	// 1. Create custom crdClientSet
	// here restConfig is your .kube/config file
	crdClientSet, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	log.Info("listing CRDs...")
	// 2. List all CRDs
	crdList, err = crdClientSet.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
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
