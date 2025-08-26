package main

import (
	"os"

	"github.com/andreis3/auth-ms/internal/auth/infra/config"
	"github.com/andreis3/auth-ms/internal/auth/infra/logger"
	"github.com/andreis3/auth-ms/internal/auth/infra/server/http"
	"github.com/andreis3/auth-ms/internal/auth/util"
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
