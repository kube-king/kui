package global

import (
	"go.uber.org/zap"
	"kube-invention/pkg/utils/logger"

	"log"
)

func init() {
	var err error
	// init zap log driver
	Log, err = logger.Init(logger.NewConfig())
	if err != nil {
		log.Printf("init zap log driver error: %s\n", err)
		return
	}
}

var (
	Log *zap.Logger
)
