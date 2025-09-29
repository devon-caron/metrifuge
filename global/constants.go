package global

import "os"

// defaults to name the variables if the env vars are not set
var (
	DEFAULT_METRICEXPORTER_CRD_NAME = "MetricExporter"
	DEFAULT_LOGEXPORTER_CRD_NAME    = "LogExporter"
	DEFAULT_LOGSOURCE_CRD_NAME      = "LogSource"
	DEFAULT_RULESET_CRD_NAME        = "RuleSet"
	DEFAULT_LOG_LEVEL               = "debug"
	DEFAULT_LOG_REPORTCALLER_STATUS = "true"
	DEFAULT_RUNNING_IN_K8S          = "true"
)

var (
	METRICEXPORTER_CRD_NAME = DEFAULT_METRICEXPORTER_CRD_NAME
	LOGEXPORTER_CRD_NAME    = DEFAULT_LOGEXPORTER_CRD_NAME
	LOGSOURCE_CRD_NAME      = DEFAULT_LOGSOURCE_CRD_NAME
	RULESET_CRD_NAME        = DEFAULT_RULESET_CRD_NAME
	LOG_LEVEL               = DEFAULT_LOG_LEVEL
	LOG_REPORTCALLER_STATUS = DEFAULT_LOG_REPORTCALLER_STATUS
	RUNNING_IN_K8S          = DEFAULT_RUNNING_IN_K8S
)

func InitConfig() {
	maybeMetricExporterCRDName := os.Getenv("MF_METRICEXPORTER_CRD_NAME")
	if maybeMetricExporterCRDName != "" {
		METRICEXPORTER_CRD_NAME = maybeMetricExporterCRDName
	}
	maybeLogExporterCRDName := os.Getenv("MF_LOGEXPORTER_CRD_NAME")
	if maybeLogExporterCRDName != "" {
		LOGEXPORTER_CRD_NAME = maybeLogExporterCRDName
	}
	maybeLogSourceCRDName := os.Getenv("MF_LOGSOURCE_CRD_NAME")
	if maybeLogSourceCRDName != "" {
		LOGSOURCE_CRD_NAME = maybeLogSourceCRDName
	}
	maybeRuleSetCRDName := os.Getenv("MF_RULESET_CRD_NAME")
	if maybeRuleSetCRDName != "" {
		RULESET_CRD_NAME = maybeRuleSetCRDName
	}
	maybeLogLevel := os.Getenv("MF_LOG_LEVEL")
	if maybeLogLevel != "" {
		LOG_LEVEL = maybeLogLevel
	}
	maybeLogReportCallerStatus := os.Getenv("MF_LOG_REPORTCALLER_STATUS")
	if maybeLogReportCallerStatus != "" {
		LOG_REPORTCALLER_STATUS = maybeLogReportCallerStatus
	}
	maybeRunningInK8s := os.Getenv("MF_RUNNING_IN_K8S")
	if maybeRunningInK8s != "" {
		RUNNING_IN_K8S = maybeRunningInK8s
	}
}
