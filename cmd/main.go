package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"task/internal/app"
	"task/internal/config"
	"task/internal/logging"
)

func main() {

	cfg := config.GetConfig()
	log := logging.GetLogger(cfg.Env, cfg.LogLevel)

	application := app.NewApp(cfg, log)
	go func() {
		application.Run()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop(context.TODO())
}
