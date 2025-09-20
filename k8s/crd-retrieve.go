package k8s

// for listing CRD, go provides client which is different from
// "kubernetes.clientset"
// This clientset will be used to list down the existing CRD
import (
	"context"
	"fmt"
	"github.com/devon-caron/metrifuge/k8s/api"
	le "github.com/devon-caron/metrifuge/k8s/api/log_exporter"
	"github.com/devon-caron/metrifuge/k8s/api/pipe"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/devon-caron/metrifuge/logger"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	log     = logger.Get()
	crdList *v1.CustomResourceDefinitionList
)

func GetK8sResources[Resource api.MetrifugeK8sResource](restConfig *rest.Config, kind, version, kindPlural string) ([]*Resource, error) {
	log.Infof("getting %v from %s", kind, restConfig.Host)
	if err := validateResources(restConfig); err != nil {
		log.Warnf("failed to validate CRDs in cluster: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
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

		resource := getResource[Resource](crdResource, kind, spec)

		resources = append(resources, resource)
	}
	return resources, nil
}

func getResource[Resource api.MetrifugeK8sResource](crdResource unstructured.Unstructured, kind string, spec map[string]interface{}) *Resource {
	var mfK8sCrdNames = []string{"RuleSet", "Pipe", "LogExporter", "MetricExporter"}

	var resource any

	switch kind {
	case mfK8sCrdNames[0]:
		resource = getRuleSet(crdResource, spec)
	case mfK8sCrdNames[1]:
		resource = getPipe(crdResource, spec)
	case mfK8sCrdNames[2]:
		resource = getLogExporter(crdResource, spec)
	case mfK8sCrdNames[3]:
		resource = getMetricExporter(crdResource, spec)
	}
	r := resource.(Resource)
	return &r
}

func getMetricExporter(crdLogExporter unstructured.Unstructured, spec map[string]any) *le.LogExporter {
	panic("getMetricExporter is not implemented yet")
}

func getLogExporter(crdLogExporter unstructured.Unstructured, spec map[string]any) *le.LogExporter {
	panic("getLogExporter is not implemented yet")
}

//func GetPipes(restConfig *rest.Config) ([]*pipe.Pipe, error) {
//	log.Infof("getting pipes from %s", restConfig.Host)
//	if err := validateResources(restConfig); err != nil {
//		log.Warnf("failed to validate CRDs in cluster: %v", err)
//	}
//
//	dynamicClient, err := dynamic.NewForConfig(restConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	gvr := schema.GroupVersionResource{
//		Group:    "metrifuge.com/k8s",
//		Version:  "v1alpha1",
//		Resource: "pipes",
//	}
//
//	crdPipeList, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	var pipes []*pipe.Pipe
//	for _, crdPipe := range crdPipeList.Items {
//		// Extract spec as map[string]any
//		spec, found, err := unstructured.NestedMap(crdPipe.Object, "spec")
//		if err != nil {
//			fmt.Printf("  Error getting spec: %v\n", err)
//			continue
//		}
//		if !found {
//			fmt.Printf("  No spec found\n")
//			continue
//		}
//
//		pipes = append(pipes, getPipe(crdPipe, spec))
//	}
//	return pipes, nil
//}

func getPipe(crdPipe unstructured.Unstructured, spec map[string]any) *pipe.Pipe {
	var matchLabelsMap map[string]any
	selectorMap, ok := spec["selector"].(map[string]any)
	if !ok {
		selectorMap = nil
		log.Debugf("  No selector found\n")
	} else {
		matchLabelsMap, ok = selectorMap["matchLabels"].(map[string]any)
		if !ok {
			log.Debugf("  No selector matchLabels\n")
		}
	}

	myPipe := &pipe.Pipe{
		APIVersion: crdPipe.GetAPIVersion(),
		Kind:       crdPipe.GetKind(),
		Metadata: api.Metadata{
			Name:      crdPipe.GetName(),
			Namespace: crdPipe.GetNamespace(),
			Labels:    crdPipe.GetLabels(),
		},
		Spec: pipe.PipeSpec{
			Selector: &api.Selector{
				MatchLabels: matchLabelsMap,
			},
			Source:   getSource(spec),
			RuleRefs: getRuleRefs(spec),
		},
	}
	return myPipe
}

func getRuleRefs( /*specMap*/ _ map[string]any) []pipe.RuleRef {
	panic("getRuleRefs function not implemented")
	return nil
}

func getSource( /*specMap*/ _ map[string]any) *pipe.Source {
	panic("getSource function (possibly along with the pipe.Source struct) not implemented")
	return nil
}

//func GetRuleSets(restConfig *rest.Config) ([]*ruleset.RuleSet, error) {
//	log.Infof("getting rules from %s", restConfig.Host)
//	if err := validateResources(restConfig); err != nil {
//		log.Warnf("failed to validate CRDs in cluster: %v", err)
//	}
//
//	dynamicClient, err := dynamic.NewForConfig(restConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	gvr := schema.GroupVersionResource{
//		Group:    "metrifuge.com/k8s",
//		Version:  "v1alpha1",
//		Resource: "rulesets",
//	}
//
//	crdRuleList, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	var ruleSets []*ruleset.RuleSet
//	for _, crdRuleSet := range crdRuleList.Items {
//		// Extract spec as map[string]any
//		spec, found, err := unstructured.NestedMap(crdRuleSet.Object, "spec")
//		if err != nil {
//			fmt.Printf("  Error getting spec: %v\n", err)
//			continue
//		}
//		if !found {
//			fmt.Printf("  No spec found\n")
//			continue
//		}
//
//		ruleSets = append(ruleSets, getRuleSet(crdRuleSet, spec))
//	}
//	return ruleSets, nil
//}

func getRuleSet(crdRuleSet unstructured.Unstructured, spec map[string]any) *ruleset.RuleSet {
	var matchLabelsMap map[string]any
	selectorMap, ok := spec["selector"].(map[string]any)
	if !ok {
		selectorMap = nil
		log.Debugf("  No selector found\n")
	} else {
		matchLabelsMap, ok = selectorMap["matchLabels"].(map[string]any)
		if !ok {
			log.Debugf("  No selector matchLabels\n")
		}
	}

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
			Selector: &api.Selector{
				MatchLabels: matchLabelsMap,
			},
			Rules: myRules,
		},
	}
	return myRuleSet
}

func getRules(_ map[string]any) []*ruleset.Rule {
	panic("getRules function not implemented")
	return nil
}

func validateResources(restConfig *rest.Config) error {

	var requiredCrdNames = []string{"RuleSet", "Pipe", "LogExporter", "MetricExporter"}

	// 1. Create custom crdClientSet
	// here restConfig is your .kube/config file
	crdClientSet, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

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
