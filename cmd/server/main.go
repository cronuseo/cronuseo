package main

import (
	"cronuseo/internal/config"
	"cronuseo/pkg/log"
	"flag"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	_, err = sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
