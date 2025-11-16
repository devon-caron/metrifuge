package log_processor

import (
	"sync"

	"github.com/devon-caron/metrifuge/k8s/api"
	logsource "github.com/devon-caron/metrifuge/k8s/api/log_source"
	"github.com/devon-caron/metrifuge/k8s/api/ruleset"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
)

type LogProcessor struct {
	sourceSets []*SourceSet
	initOnce   sync.Once
	log        *logrus.Logger
}

type SourceSet struct {
	source api.Source
	rules  []*api.Rule
}

func (lp *LogProcessor) Initialize(logSources []*logsource.LogSource, ruleSets []*ruleset.RuleSet, log *logrus.Logger) {
	if logSources == nil || ruleSets == nil || log == nil {
		logrus.Fatalf("log processor initialization failed: logSources: %v, ruleSets: %v, log: %v", logSources, ruleSets, log)
	}

	lp.log = log

	lp.initOnce.Do(func() {
		lp.sourceSets = make([]*SourceSet, 0)

		for _, rs := range ruleSets {
			lp.log.Infof("processing rule set: %v", rs)

			lp.log.Debugf("spec: %v", rs.Spec)
			lp.log.Debugf("selector: %v", rs.Spec.Selector)
			lp.log.Debugf("matchlabels: %v", rs.Spec.Selector.MatchLabels)

			selectorLabels := rs.Spec.Selector.MatchLabels

			selector := labels.Set(selectorLabels).AsSelector()
			for _, ls := range logSources {
				lp.log.Infof("processing log source: %v", ls)
				sourceLabels := ls.Metadata.Labels
				lp.log.Infof("source labels: %v", sourceLabels)
				if selector.Matches(labels.Set(sourceLabels)) {
					set := &SourceSet{
						rules: rs.Spec.Rules,
					}

					switch ls.Spec.Source.Type {
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
	})
}
