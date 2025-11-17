package global

import "os"

// defaults to name the variables if the env vars are not set
var (
	DEFAULT_LOGSOURCE_CRD_NAME      = "LogSource"
	DEFAULT_RULESET_CRD_NAME        = "RuleSet"
	DEFAULT_LOG_LEVEL               = "debug"
	DEFAULT_LOG_REPORTCALLER_STATUS = "true"
	DEFAULT_RUNNING_IN_K8S          = "true"
	DEFAULT_EXPORTER_CRD_NAME       = "Exporter"
	DEFAULT_REFRESH_INTERVAL        = "60"
	DEFAULT_LOG_SOURCE_RETRIES      = "5"
	DEFAULT_LOG_SOURCE_DELAY        = "5"
)

var (
	LOGSOURCE_CRD_NAME      = DEFAULT_LOGSOURCE_CRD_NAME
	RULESET_CRD_NAME        = DEFAULT_RULESET_CRD_NAME
	LOG_LEVEL               = DEFAULT_LOG_LEVEL
	LOG_REPORTCALLER_STATUS = DEFAULT_LOG_REPORTCALLER_STATUS
	RUNNING_IN_K8S          = DEFAULT_RUNNING_IN_K8S
	EXPORTER_CRD_NAME       = DEFAULT_EXPORTER_CRD_NAME
	REFRESH_INTERVAL        = DEFAULT_REFRESH_INTERVAL
	LOG_SOURCE_RETRIES      = DEFAULT_LOG_SOURCE_RETRIES
	LOG_SOURCE_DELAY        = DEFAULT_LOG_SOURCE_DELAY
)

func InitConfig() {
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
	maybeRefreshInterval := os.Getenv("MF_REFRESH_INTERVAL")
	if maybeRefreshInterval != "" {
		REFRESH_INTERVAL = maybeRefreshInterval
	}
	maybeExporterCRDName := os.Getenv("MF_EXPORTER_CRD_NAME")
	if maybeExporterCRDName != "" {
		EXPORTER_CRD_NAME = maybeExporterCRDName
	}
	maybeLogSourceRetries := os.Getenv("MF_LOG_SOURCE_RETRIES")
	if maybeLogSourceRetries != "" {
		LOG_SOURCE_RETRIES = maybeLogSourceRetries
	}
	maybeLogSourceDelay := os.Getenv("MF_LOG_SOURCE_DELAY")
	if maybeLogSourceDelay != "" {
		LOG_SOURCE_DELAY = maybeLogSourceDelay
	}
}
