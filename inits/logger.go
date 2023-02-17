package inits

import (
	"email2misskey/config"
	"email2misskey/global"
	"fmt"
	"go.uber.org/zap"
)

func Logger() error {

	var err error

	var logger *zap.Logger

	// Prepare logger
	if config.Config.System.Production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Flush logs
	defer logger.Sync() // Unable to handle errors here

	// Sugar it
	global.Logger = logger.Sugar()

	return nil
}
