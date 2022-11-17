package main

import (
	"cronuseo/internal/config"
	"flag"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	// load application configurations
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		os.Exit(-1)
	}

	//connect db
	_, err = sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		os.Exit(-1)
	}

}
