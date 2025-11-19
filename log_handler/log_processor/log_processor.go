package log_processor

import (
	"fmt"

	"github.com/devon-caron/metrifuge/k8s/api"
	logsource "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
)

type LogProcessor struct {
	sourceSets []*LogSet
	log        *logrus.Logger
}

type LogSet struct {
	source api.Source
	rules  []*api.Rule
}

type ProcessedData struct {
	ForwardLog string
	Metric     *MetricData
}

type MetricData struct {
	Name       string
	Value      string
	Template   *api.MetricTemplate
	Attributes []api.Attribute
}

func (lp *LogProcessor) Initialize(logSources []*logsource.LogSource, ruleSets []*ruleset.RuleSet, log *logrus.Logger) {
	if logSources == nil || ruleSets == nil || log == nil {
		logrus.Fatalf("log processor initialization failed: logSources: %v, ruleSets: %v, log: %v", logSources, ruleSets, log)
	}

	lp.log = log

	lp.sourceSets = make([]*LogSet, 0)

	lp.Update(logSources, ruleSets)
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
				set := &LogSet{
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

func (lp *LogProcessor) FindLogSet(source api.Source) (*LogSet, error) {
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

func (ls *LogSet) ProcessLogs(logs []string) []ProcessedData {
	processedData := make([]ProcessedData, 0)
	for _, log := range logs {
		for _, rule := range ls.rules {
			processedDataPoint, err := processLog(log, rule)
			if err != nil {
				logrus.Errorf("failed to process log: %v", err)
				continue
			}
			processedData = append(processedData, processedDataPoint)
		}
	}
	return processedData
}

// TODO needs implementation
func processLog(log string, rule *api.Rule) (ProcessedData, error) {
	// logrus.Infof("fake processing log: %v", log)
	// logrus.Infof("fake processing rule: %v", rule)
	return ProcessedData{}, nil
}
