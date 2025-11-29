package log_processor

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/devon-caron/metrifuge/k8s/api"
	logsource "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/grok"
	"go.opentelemetry.io/otel/attribute"
	"k8s.io/apimachinery/pkg/labels"
)

type LogProcessor struct {
	sourceSets []*SourceRuleUnion
	log        *logrus.Logger
	g          *grok.Grok
}

type SourceRuleUnion struct {
	source api.Source
	rules  []*api.Rule
}

func (lp *LogProcessor) Initialize(logSources []logsource.LogSource, ruleSets []ruleset.RuleSet, log *logrus.Logger) {
	if logSources == nil {
		logrus.Fatalf("log processor initialization failed, logSources triggered nil: logSources: %v, ruleSets: %v, log: %v", logSources, ruleSets, log)
	}

	if ruleSets == nil {
		logrus.Fatalf("log processor initialization failed, ruleSets triggered nil: logSources: %v, ruleSets: %v, log: %v", logSources, ruleSets, log)
	}

	if log == nil {
		logrus.Fatalf("log processor initialization failed, log triggered nil: logSources: %v, ruleSets: %v, log: %v", logSources, ruleSets, log)
	}

	lp.log = log

	lp.sourceSets = make([]*SourceRuleUnion, 0)

	lp.Update(logSources, ruleSets)

	if g, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true}); err != nil {
		logrus.Fatalf("failed to initialize grok: %v", err)
	} else {
		lp.g = g
	}
}

func (lp *LogProcessor) Update(logSources []logsource.LogSource, ruleSets []ruleset.RuleSet) {
	for _, rs := range ruleSets {
		lp.log.Infof("processing rule set: %v", rs)

		lp.log.Debugf("spec: %v", rs.Spec)
		lp.log.Debugf("selector: %v", rs.Spec.Selector)
		lp.log.Debugf("matchlabels: %v", rs.Spec.Selector.MatchLabels)

		selectorLabels := rs.Spec.Selector.MatchLabels

		selector := labels.Set(selectorLabels).AsSelector()
		for _, ls := range logSources {
			lp.log.Infof("processing log source: %+v", ls)
			sourceLabels := ls.Metadata.Labels
			lp.log.Infof("source labels: %v", sourceLabels)
			lp.log.Infof("type: %v", ls.Spec.Type)
			if selector.Matches(labels.Set(sourceLabels)) {
				set := &SourceRuleUnion{
					rules: rs.Spec.Rules,
				}

				switch ls.Spec.Type {
				case "PodSource":
					set.source = ls.Spec.Source.PodSource
				case "PVCSource":
					set.source = ls.Spec.Source.PVCSource
				case "LocalSource":
					set.source = ls.Spec.Source.LocalSource
				case "CmdSource":
					set.source = ls.Spec.Source.CmdSource
				default:
					lp.log.Errorf("unknown log source type: %s", ls.Spec.Source.Type)
					return
				}

				lp.sourceSets = append(lp.sourceSets, set)
				lp.log.Infof("added source set: %v", set)
			}
		}
	}
}

func (lp *LogProcessor) FindSRU(source api.Source) (*SourceRuleUnion, error) {
	for i, set := range lp.sourceSets {
		// TODO this is a costly operation, needs improvement
		lp.log.Infof("checking source #%v: %v", i+1, set.source.GetSourceInfo())
		lp.log.Infof("against desired source: %v", source.GetSourceInfo())
		if set.source.GetSourceInfo() == source.GetSourceInfo() {
			return set, nil
		}
	}
	return nil, fmt.Errorf("log set not found for source: %v", source.GetSourceInfo())
}

func (lp *LogProcessor) ProcessLogsWithSRU(sru *SourceRuleUnion, logs []string, lsName string, lsNamespace string) []api.ProcessedDataItem {
	totalProcessedDataItems := make([]api.ProcessedDataItem, 0)
	baseCtx := context.WithValue(context.Background(), "name", lsName)
	baseCtx = context.WithValue(baseCtx, "namespace", lsNamespace)
	for _, log := range logs {
		for _, rule := range sru.rules {
			processedDataItems, err := lp.processLog(baseCtx, log, rule)
			if err != nil {
				lp.log.Errorf("failed to process log: %v", err)
				continue
			}
			totalProcessedDataItems = append(totalProcessedDataItems, processedDataItems...)
		}
	}

	lp.log.Debugf("processed %d logs into %d data items", len(logs), len(totalProcessedDataItems))
	if len(totalProcessedDataItems) > 0 {
		lp.log.Debugf("first data item: %+v", totalProcessedDataItems[0])
	}
	return totalProcessedDataItems
}

// TODO needs implementation
func (lp *LogProcessor) processLog(ctx context.Context, logMsg string, rule *api.Rule) ([]api.ProcessedDataItem, error) {

	var srcInfo = api.LogSourceInfo{}

	lsName, ok := ctx.Value("name").(string)
	if !ok {
		return []api.ProcessedDataItem{}, fmt.Errorf("missing name in context")
	}
	lsNamespace, ok := ctx.Value("namespace").(string)
	if !ok {
		return []api.ProcessedDataItem{}, fmt.Errorf("missing namespace in context")
	}
	srcInfo.Name = lsName
	srcInfo.Namespace = lsNamespace

	values, err := lp.g.Parse(rule.Pattern, logMsg)
	if err != nil {
		return []api.ProcessedDataItem{}, err
	}

	// debug
	if rand.IntN(10) == 0 {
		lp.log.Debugf("parsed log: %v", logMsg)
		lp.log.Debugf("pattern: %v", rule.Pattern)
		for k, v := range values {
			lp.log.Debugf("%v: %v", k, v)
		}
	}

	metricData, err := lp.createMetricData(values, rule.Metrics)
	if err != nil {
		return []api.ProcessedDataItem{}, err
	}

	lp.log.Debugf("created %d metric data items", len(metricData))

	processedLogMsg := ""
	processedDataItems := make([]api.ProcessedDataItem, 0)
	switch strings.ToLower(rule.Action) {
	case "forward":
		processedLogMsg = logMsg
	case "discard":
		// processedLogMsg = ""
		lp.log.Debugf("Discard Action No-Op")
	case "conditional":
		processedLogMsg, processedDataItems, err = lp.processConditional(ctx, logMsg, values, rule, rule.Conditional)
		if err != nil {
			return []api.ProcessedDataItem{}, err
		}
	default:
		return []api.ProcessedDataItem{}, fmt.Errorf("unknown action: %v", rule.Action)
	}

	for _, metric := range metricData {
		processedDataItems = append(processedDataItems, api.ProcessedDataItem{
			ForwardLog:    processedLogMsg,
			Metric:        metric,
			LogSourceInfo: srcInfo,
		})
	}

	return processedDataItems, nil
}

func (lp *LogProcessor) createMetricData(values map[string]string, metrics []api.MetricTemplate) ([]*api.MetricData, error) {

	myMetricDataList := make([]*api.MetricData, 0)

	for _, metricTemplate := range metrics {
		lp.log.Debugf("processing metric template: %s", metricTemplate.Name)
		lp.log.Debugf("metric template details: %+v", metricTemplate)
		metricData := &api.MetricData{
			Name:       metricTemplate.Name,
			Kind:       metricTemplate.Kind,
			Attributes: make([]attribute.KeyValue, 0),
		}
		var err error
		// TODO improve shitty implementation
		switch strings.ToLower(metricTemplate.Value.Type) {
		case "int64":
			lp.log.Debugf("processing int64 metric: %s", metricTemplate.Name)
			if metricTemplate.Value.GrokKey == "" {
				metricData.ValueInt, err = strconv.ParseInt(metricTemplate.Value.ManualValue, 10, 64)
			} else {
				metricData.ValueInt, err = strconv.ParseInt(values[metricTemplate.Value.GrokKey], 10, 64)
			}
			if err != nil {
				return []*api.MetricData{}, fmt.Errorf("failed to parse int64 metric value: %v", err)
			}
		case "float64":
			lp.log.Debugf("processing float64 metric: %s", metricTemplate.Name)
			if metricTemplate.Value.GrokKey == "" {
				metricData.ValueFloat, err = strconv.ParseFloat(metricTemplate.Value.ManualValue, 64)
			} else {
				metricData.ValueFloat, err = strconv.ParseFloat(values[metricTemplate.Value.GrokKey], 64)
			}
			if err != nil {
				return []*api.MetricData{}, fmt.Errorf("failed to parse float64 metric value: %v", err)
			}
		default:
			return []*api.MetricData{}, fmt.Errorf("unknown metric value type: %v", metricTemplate.Value.Type)
		}

		for _, currAttribute := range metricTemplate.Attributes {
			switch strings.ToLower(currAttribute.Value.Type) {
			case "int64":
				lp.log.Debugf("processing int64 attribute: %s", currAttribute.Key)
				if currAttribute.Value.GrokKey == "" {
					attrValue, parseErr := strconv.ParseInt(currAttribute.Value.ManualValue, 10, 64)
					if parseErr != nil {
						return []*api.MetricData{}, fmt.Errorf("failed to parse int64 attribute value: %v", parseErr)
					}
					metricData.Attributes = append(metricData.Attributes, attribute.Int64(currAttribute.Key, attrValue))
				} else {
					attrValue, parseErr := strconv.ParseInt(values[currAttribute.Value.GrokKey], 10, 64)
					if parseErr != nil {
						return []*api.MetricData{}, fmt.Errorf("failed to parse int64 attribute value: %v", parseErr)
					}
					metricData.Attributes = append(metricData.Attributes, attribute.Int64(currAttribute.Key, int64(attrValue)))
				}
				if err != nil {
					return []*api.MetricData{}, fmt.Errorf("failed to parse int64 metric attribute value: %v", err)
				}
			case "float64":
				lp.log.Debugf("processing float64 attribute: %s", currAttribute.Key)
				if currAttribute.Value.GrokKey == "" {
					attrValue, parseErr := strconv.ParseFloat(currAttribute.Value.ManualValue, 64)
					if parseErr != nil {
						return []*api.MetricData{}, fmt.Errorf("failed to parse float64 attribute value: %v", parseErr)
					}
					metricData.Attributes = append(metricData.Attributes, attribute.Float64(currAttribute.Key, attrValue))
				} else {
					attrValue, parseErr := strconv.ParseFloat(values[currAttribute.Value.GrokKey], 64)
					if parseErr != nil {
						return []*api.MetricData{}, fmt.Errorf("failed to parse float64 attribute value: %v", parseErr)
					}
					metricData.Attributes = append(metricData.Attributes, attribute.Float64(currAttribute.Key, attrValue))
				}
				if err != nil {
					return []*api.MetricData{}, fmt.Errorf("failed to parse float64 metric attribute value: %v", err)
				}
			case "string":
				lp.log.Debugf("processing string attribute: %s", currAttribute.Key)
				if currAttribute.Value.GrokKey == "" {
					attrValue := currAttribute.Value.ManualValue
					metricData.Attributes = append(metricData.Attributes, attribute.String(currAttribute.Key, attrValue))
				} else {
					attrValue := values[currAttribute.Value.GrokKey]
					metricData.Attributes = append(metricData.Attributes, attribute.String(currAttribute.Key, attrValue))
				}
			default:
				return []*api.MetricData{}, fmt.Errorf("unknown metric attribute value type: %v", currAttribute.Value.Type)
			}
		}

		myMetricDataList = append(myMetricDataList, metricData)
		// if rand.IntN(20) == 0 {
		lp.log.Debugf("metric data: %+v", metricData)
		// }
	}
	return myMetricDataList, nil
}

func (lp *LogProcessor) processConditional(ctx context.Context, logMsg string, values map[string]string, rule *api.Rule, conditional *api.Conditional) (string, []api.ProcessedDataItem, error) {

	lp.log.Debugf("Evaluating conditional rule with pattern %s with field1: %v, operator: %s", rule.Pattern, conditional.Field1, conditional.Operator)
	lp.log.Debugf("conditional: %+v", conditional)

	var srcInfo = api.LogSourceInfo{}

	lsName, ok := ctx.Value("name").(string)
	if !ok {
		return "", []api.ProcessedDataItem{}, fmt.Errorf("missing name in context")
	}
	lsNamespace, ok := ctx.Value("namespace").(string)
	if !ok {
		return "", []api.ProcessedDataItem{}, fmt.Errorf("missing namespace in context")
	}
	srcInfo.Name = lsName
	srcInfo.Namespace = lsNamespace

	var f1Str, f2Str string

	// Extract string values from FieldValue structs
	if conditional.Field1.GrokKey != "" {
		f1Str = values[conditional.Field1.GrokKey]
	} else {
		f1Str = conditional.Field1.ManualValue
	}

	if conditional.Field2.GrokKey != "" {
		f2Str = values[conditional.Field2.GrokKey]
	} else {
		f2Str = conditional.Field2.ManualValue
	}

	op := conditional.Operator

	lp.log.Debugf("validating conditional: field1='%s', field2='%s', operator='%s'", f1Str, f2Str, op)

	var err error
	if err = lp.validateFields(f1Str, f2Str, op); err != nil {
		return "", nil, fmt.Errorf("conditional validation failed: %w", err)
	}

	lp.log.Debugf("fields validation passed")

	lp.log.Debugf("evaluating conditional result")

	var result bool
	if result, err = evaluateConditional(f1Str, f2Str, op); err != nil {
		return "", nil, fmt.Errorf("conditional evaluation failed: %w", err)
	}

	lp.log.Debugf("conditional result: %t", result)

	var selectedAction string
	if result {
		selectedAction = conditional.ActionTrue
	} else {
		selectedAction = conditional.ActionFalse
	}

	var fwdLog string
	var extraDataItems []api.ProcessedDataItem
	switch strings.ToLower(selectedAction) {
	case "forward":
		fwdLog = logMsg
	case "discard":
		fwdLog = ""
	case "conditional":
		var resultConditional *api.Conditional
		if result {
			resultConditional = conditional.ConditionalTrue
		} else {
			resultConditional = conditional.ConditionalFalse
		}
		fwdLog, extraDataItems, err = lp.processConditional(ctx, logMsg, values, rule, resultConditional)
		if err != nil {
			return "", nil, fmt.Errorf("nested conditional processing failed: %w", err)
		}
	default:
		return "", nil, fmt.Errorf("unknown action: %s", selectedAction)
	}

	lp.log.Debugf("selected action: %s, result: %t", selectedAction, result)

	var resultMetrics []api.MetricTemplate
	if result {
		resultMetrics = conditional.MetricsTrue
	} else {
		resultMetrics = conditional.MetricsFalse
	}

	lp.log.Debugf("selected metrics for result %t: %v", result, resultMetrics)
	metricData, err := lp.createMetricData(values, resultMetrics)
	if err != nil {
		return "", nil, fmt.Errorf("metric data creation failed: %w", err)
	}

	lp.log.Debugf("created %d metric data items", len(metricData))

	var processedDataItems = make([]api.ProcessedDataItem, 0)
	for _, metric := range metricData {
		processedDataItems = append(processedDataItems, api.ProcessedDataItem{
			ForwardLog:    fwdLog,
			Metric:        metric,
			LogSourceInfo: srcInfo,
		})
	}

	processedDataItems = append(processedDataItems, extraDataItems...)

	return fwdLog, processedDataItems, nil
}

func (lp *LogProcessor) validateFields(field1, field2, op string) error {
	// Check that field1 and field2 are valid based on the operator
	switch op {
	case "Equals", "DoesNotEqual":
		// These operators require both fields to be present and comparable
		if field1 == "" {
			return fmt.Errorf("field1 is required for operator %s", op)
		}
		if field2 == "" {
			return fmt.Errorf("field2 is required for operator %s", op)
		}
	case "Exists", "DoesNotExist":
		lp.log.Debugf("Exists/DoesNotExist operators don't require field validation")
	case "LessThan", "GreaterThan", "GreaterThanOrEqualTo", "LessThanOrEqualTo":
		// These operators require both fields to be present, parseable as integers, and comparable
		if field1 == "" {
			return fmt.Errorf("field1 is required for operator %s", op)
		}
		if field2 == "" {
			return fmt.Errorf("field2 is required for operator %s", op)
		}
		// Try to parse as integers for comparison
		_, err1 := strconv.Atoi(field1)
		_, err2 := strconv.Atoi(field2)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("field1 and field2 must be parseable as integers for operator %s", op)
		}
	default:
		return fmt.Errorf("unsupported operator: %s", op)
	}

	return nil
}

func evaluateConditional(f1Str, f2Str, op string) (bool, error) {
	switch op {
	case "Equals":
		return f1Str == f2Str, nil
	case "DoesNotEqual":
		return f1Str != f2Str, nil
	case "Exists":
		return f1Str != "", nil
	case "DoesNotExist":
		return f1Str == "", nil
	case "LessThan":
		i1, err1 := strconv.Atoi(f1Str)
		i2, err2 := strconv.Atoi(f2Str)
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("cannot compare non-numeric values: %s and %s", f1Str, f2Str))
		}
		return i1 < i2, nil
	case "GreaterThan":
		i1, err1 := strconv.Atoi(f1Str)
		i2, err2 := strconv.Atoi(f2Str)
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("cannot compare non-numeric values: %s and %s", f1Str, f2Str))
		}
		return i1 > i2, nil
	case "GreaterThanOrEqualTo":
		i1, err1 := strconv.Atoi(f1Str)
		i2, err2 := strconv.Atoi(f2Str)
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("cannot compare non-numeric values: %s and %s", f1Str, f2Str))
		}
		return i1 >= i2, nil
	case "LessThanOrEqualTo":
		i1, err1 := strconv.Atoi(f1Str)
		i2, err2 := strconv.Atoi(f2Str)
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("cannot compare non-numeric values: %s and %s", f1Str, f2Str))
		}
		return i1 <= i2, nil
	default:
		panic(fmt.Errorf("unsupported operator: %s", op))
	}
}
