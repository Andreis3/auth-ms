package main

import (
	"os"

	"github.com/andreis3/auth-ms/internal/infra/config"
	"github.com/andreis3/auth-ms/internal/infra/logger"
	"github.com/andreis3/auth-ms/internal/infra/server/http"
	"github.com/andreis3/auth-ms/internal/util"
)

func main() {
	conf := config.LoadConfig()
	log := logger.NewLogger()

	if conf == nil {
		log.CriticalText("Failed to load configuration")
		os.Exit(util.ExitFailure)
	}

	serverWeb := http.NewServer(conf, *log)

	go serverWeb.Start()

	serverWeb.GracefulShutdown()
}
