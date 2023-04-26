package test

import (
	"log"
	"os"
	"testing"

	"github.com/shashimalcse/cronuseo/internal/config"
	db "github.com/shashimalcse/cronuseo/internal/db/mongo"
	"github.com/shashimalcse/cronuseo/internal/logger"
)

func DB(t *testing.T) *db.MongoDB {

	cfg, err := config.Load("../../config/local-debug.yml")
	if err != nil {
		log.Fatal("Error while loading config for test.")
		os.Exit(-1)
	}
	logger := logger.Init(cfg)
	mongo := db.Init(cfg, logger)

	return mongo

}
