package inits

import (
	"email2misskey/config"
	"email2misskey/consts"
	"email2misskey/global"
	"email2misskey/handlers"
	"fmt"
	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
)

func SMTP() error {
	// Config hosts
	cfg := &guerrilla.AppConfig{
		AllowedHosts: config.Config.EMail.Host,
	}
	// Config logger
	if config.Config.System.Production {
		cfg.LogLevel = "error"
	} else {
		cfg.LogLevel = "debug"
	}
	// Config server
	sc := guerrilla.ServerConfig{
		ListenInterface: config.Config.System.Listen,
		IsEnabled:       true,
	}
	cfg.Servers = append(cfg.Servers, sc)
	// Config backend
	additionalSaveProcess := ""
	if !config.Config.System.Production {
		additionalSaveProcess = "|Debugger"
	}
	bcfg := backends.BackendConfig{
		"save_process":       fmt.Sprintf("HeadersParser|Header|Hasher%s|%s", additionalSaveProcess, consts.ProcessorID),
		"log_received_mails": !config.Config.System.Production,
	}
	cfg.BackendConfig = bcfg
	d := guerrilla.Daemon{
		Config: cfg,
	}

	// Add handlers
	d.AddProcessor(consts.ProcessorID, handlers.IncomingEMAil())

	// Start listening
	err := d.Start()
	if err != nil {
		global.Logger.Errorf("Failed to start SMTP server with error: %v", err)
	} else {
		global.Logger.Infof("SMTP server started!")
	}

	global.GDaemon = &d

	return err
}
