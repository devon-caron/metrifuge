package core

import (
	"github.com/devon-caron/metrifuge/api"
	log "github.com/sirupsen/logrus"
)

func Start() {
	log.Info("starting api")
	api.StartApi()

}
