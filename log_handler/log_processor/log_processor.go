package log_processor

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/devon-caron/metrifuge/k8s/api"
	logsource "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/grok"
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

type ProcessedDataItem struct {
	ForwardLog string
	Metric     *MetricData
}

type MetricData struct {
	Name             string
	Kind             string
	ValueInt         int64
	ValueFloat       float64
	AttributesInt    map[string]int64
	AttributesFloat  map[string]float64
	AttributesString map[string]string
}

func (lp *LogProcessor) Initialize(logSources []*logsource.LogSource, ruleSets []*ruleset.RuleSet, log *logrus.Logger) {
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

func (lp *LogProcessor) Update(logSources []*logsource.LogSource, ruleSets []*ruleset.RuleSet) {

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

func (lp *LogProcessor) ProcessLogsWithSRU(sru *SourceRuleUnion, logs []string) []ProcessedDataItem {
	totalProcessedDataItems := make([]ProcessedDataItem, 0)
	for _, log := range logs {
		for _, rule := range sru.rules {
			processedDataItems, err := lp.processLog(log, rule)
			if err != nil {
				logrus.Errorf("failed to process log: %v", err)
				continue
			}
			totalProcessedDataItems = append(totalProcessedDataItems, processedDataItems...)
		}
	}
	return totalProcessedDataItems
}

// TODO needs implementation
func (lp *LogProcessor) processLog(logMsg string, rule *api.Rule) ([]ProcessedDataItem, error) {
	values, err := lp.g.Parse(rule.Pattern, logMsg)
	if err != nil {
		return []ProcessedDataItem{}, err
	}

	// debug
	if rand.IntN(10) == 0 {
		lp.log.Debugf("parsed log: %v", logMsg)
		lp.log.Debugf("pattern: %v", rule.Pattern)
		for k, v := range values {
			lp.log.Debugf("%v: %v", k, v)
		}
	}

	metricData, err := lp.createMetricData(values, rule)
	if err != nil {
		return []ProcessedDataItem{}, err
	}

	processedDataItems := make([]ProcessedDataItem, 0)
	for _, metric := range metricData {
		processedDataItems = append(processedDataItems, ProcessedDataItem{
			ForwardLog: logMsg,
			Metric:     metric,
		})
	}

	return processedDataItems, nil
}

func (lp *LogProcessor) createMetricData(values map[string]string, rule *api.Rule) ([]*MetricData, error) {

	myMetricDataList := make([]*MetricData, 0)

	for _, metricTemplate := range rule.Metrics {
		metricData := &MetricData{
			Name:             metricTemplate.Name,
			Kind:             metricTemplate.Kind,
			AttributesInt:    make(map[string]int64),
			AttributesFloat:  make(map[string]float64),
			AttributesString: make(map[string]string),
		}
		var err error
		// TODO improve shitty implementation
		switch metricTemplate.Value.Type {
		case "Int64":
			if metricTemplate.Value.GrokKey == "" {
				metricData.ValueInt, err = strconv.ParseInt(metricTemplate.Value.ManualValue, 10, 64)
			} else {
				metricData.ValueInt, err = strconv.ParseInt(values[metricTemplate.Value.GrokKey], 10, 64)
			}
			if err != nil {
				return []*MetricData{}, fmt.Errorf("failed to parse int64 metric value: %v", err)
			}
		case "Float64":
			if metricTemplate.Value.GrokKey == "" {
				metricData.ValueFloat, err = strconv.ParseFloat(metricTemplate.Value.ManualValue, 64)
			} else {
				metricData.ValueFloat, err = strconv.ParseFloat(values[metricTemplate.Value.GrokKey], 64)
			}
			if err != nil {
				return []*MetricData{}, fmt.Errorf("failed to parse float64 metric value: %v", err)
			}
		default:
			return []*MetricData{}, fmt.Errorf("unknown metric value type: %v", metricTemplate.Value.Type)
		}

		for _, attribute := range metricTemplate.Attributes {
			switch attribute.Value.Type {
			case "Int64":
				if attribute.Value.GrokKey == "" {
					metricData.AttributesInt[attribute.Key], err = strconv.ParseInt(attribute.Value.ManualValue, 10, 64)
				} else {
					metricData.AttributesInt[attribute.Key], err = strconv.ParseInt(values[attribute.Value.GrokKey], 10, 64)
				}
				if err != nil {
					return []*MetricData{}, fmt.Errorf("failed to parse int64 metric attribute value: %v", err)
				}
			case "Float64":
				if attribute.Value.GrokKey == "" {
					metricData.AttributesFloat[attribute.Key], err = strconv.ParseFloat(attribute.Value.ManualValue, 64)
				} else {
					metricData.AttributesFloat[attribute.Key], err = strconv.ParseFloat(values[attribute.Value.GrokKey], 64)
				}
				if err != nil {
					return []*MetricData{}, fmt.Errorf("failed to parse float64 metric attribute value: %v", err)
				}
			case "String":
				if attribute.Value.GrokKey == "" {
					metricData.AttributesString[attribute.Key] = attribute.Value.ManualValue
				} else {
					metricData.AttributesString[attribute.Key] = values[attribute.Value.GrokKey]
				}
			default:
				return []*MetricData{}, fmt.Errorf("unknown metric attribute value type: %v", attribute.Value.Type)
			}
		}

		myMetricDataList = append(myMetricDataList, metricData)
		if rand.IntN(20) == 0 {
			lp.log.Debugf("metric data: %+v", metricData)
		}
	}
	return myMetricDataList, nil
}
