package logger

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/devon-caron/metrifuge/global"
	"github.com/sirupsen/logrus"
)

var (
	logLevel logrus.Level
	logger   *logrus.Logger
	once     sync.Once
)

func Get() *logrus.Logger {
	once.Do(func() {
		logLevel = initLogLevel(global.LOG_LEVEL)

		logger = logrus.StandardLogger()

		logger.SetLevel(logLevel)

		LOG_REPORTCALLER_STATUS, err := strconv.ParseBool(global.LOG_REPORTCALLER_STATUS)
		if err != nil {
			log.Fatalf("Error parsing global.LOG_REPORTCALLER_STATUS: %s", err)
		}
		logger.SetReportCaller(LOG_REPORTCALLER_STATUS)

		logger.Info("logger initialized")
	})
	return logger
}

func initLogLevel(llStr string) logrus.Level {
	switch strings.ToLower(llStr) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		panic("MF_LOG_LEVEL env variable is not valid: " + llStr)
	}
}
