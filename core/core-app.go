package core

import (
	"github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/devon-caron/metrifuge/receiver"
	"github.com/sirupsen/logrus"
)

var (
	lr  *receiver.LogReceiver
	log *logrus.Logger
)

func Start() {
	log = logger.Get()
	log.Info("starting api")
	api.StartApi()
	lr = receiver.GetLogReceiver()
}
