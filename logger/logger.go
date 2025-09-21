package logger

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	logLevel                   logrus.Level
	logger                     *logrus.Logger
	once                       sync.Once
	MF_LOG_LEVEL               string
	MF_LOG_REPORTCALLER_STATUS bool
)

func Get() *logrus.Logger {
	once.Do(func() {
		MF_LOG_LEVEL = os.Getenv("MF_LOG_LEVEL")
		logLevel = initLogLevel(MF_LOG_LEVEL)

		logger = logrus.StandardLogger()

		logger.SetLevel(logLevel)

		var err error
		MF_LOG_REPORTCALLER_STATUS, err = strconv.ParseBool(os.Getenv("MF_LOG_REPORTCALLER_STATUS"))
		if err != nil {
			log.Fatalf("Error parsing MF_LOG_REPORTCALLER_STATUS: %s", err)
		}
		logger.SetReportCaller(MF_LOG_REPORTCALLER_STATUS)

		logger.Info("logger initialized")
	})
	return logger
}

func initLogLevel(llStr string) logrus.Level {
	switch llStr {
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
