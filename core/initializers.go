package core

import "github.com/devon-caron/metrifuge/k8s/crd"

func initPipes(isK8s bool) []*crd.Pipe {
	if !isK8s {
		return []*crd.Pipe{}
	}

	return []*crd.Pipe{}
}

func initRules(isK8s bool) []*crd.Rule {
	if !isK8s {
		return []*crd.Rule{}
	}

	return []*crd.Rule{}
}

func initMetricExporters(isK8s bool) []*crd.MetricExporter {
	if !isK8s {
		return []*crd.MetricExporter{}
	}

	return []*crd.MetricExporter{}
}

func initLogExporters(isK8s bool) []*crd.LogExporter {
	if !isK8s {
		return []*crd.LogExporter{}
	}

	return []*crd.LogExporter{}
}
