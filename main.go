package main

import (
	"email2misskey/global"
	"email2misskey/inits"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Start initializing

	log.Println("Start initializing...")
	if err := inits.Config(); err != nil {
		log.Fatal(err)
	}
	if err := inits.Logger(); err != nil {
		log.Fatal(err)
	}
	global.Logger.Infof("Logger initialized, switch to here")
	if err := inits.Redis(); err != nil {
		global.Logger.Fatal(err)
	}
	if err := inits.Misskey(); err != nil {
		global.Logger.Fatal(err)
	}
	if err := inits.SMTP(); err != nil {
		global.Logger.Fatal(err)
	}
	global.Logger.Infof("System initialized successfully.")

	// Block process till finish
	channelComplete := setupSignalHandling()
	<-*channelComplete
}

func setupSignalHandling() *chan struct{} {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	channelComplete := make(chan struct{})

	go func() {
		<-signalChannel
		shutdown()
		channelComplete <- struct{}{}
	}()

	return &channelComplete
}

func shutdown() {
	global.Logger.Infof("System shutting down...")
	if global.GDaemon != nil {
		global.GDaemon.Shutdown()
	}
}
