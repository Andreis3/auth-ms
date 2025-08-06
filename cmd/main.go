package main

import (
	"os"

	"github.com/andreis3/customers-ms/internal/infra/adapters/logger"
	"github.com/andreis3/customers-ms/internal/infra/configs"
	"github.com/andreis3/customers-ms/internal/infra/server/web"
	"github.com/andreis3/customers-ms/internal/util"
)

func main() {
	conf := configs.LoadConfig()
	log := logger.NewLogger()

	if conf == nil {
		log.CriticalText("Failed to load configuration")
		os.Exit(util.ExitFailure)
	}

	serverWeb := web.NewServer(conf, *log)

	go serverWeb.Start()

	serverWeb.GracefulShutdown()
}
