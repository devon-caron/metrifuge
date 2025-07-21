package core

import (
	"github.com/devon-caron/metrifuge/api"
	"github.com/devon-caron/metrifuge/receiver"
	log "github.com/sirupsen/logrus"
)

var (
	lr *receiver.LogReceiver
)

func Start() {
	log.Info("starting api")
	api.StartApi()
	lr = receiver.GetLogReceiver()
}
