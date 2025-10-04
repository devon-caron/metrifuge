package api

import (
	"github.com/devon-caron/metrifuge/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func StartApi() {
	log = logger.Get()
}
