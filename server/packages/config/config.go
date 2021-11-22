package config

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/joho/godotenv"
)

const (
	POSTGRES_USER        = "POSTGRES_USER"
	POSTGRES_PASSWORD    = "POSTGRES_PASSWORD"
	POSTGRES_DB          = "POSTGRES_DB"
	CLIENT_URL           = "CLIENT_URL"
	SERVER_PORT          = "SERVER_PORT"
	JWT_KEY              = "JWT_KEY"
	RUN_MIGRATION        = "RUN_MIGRATION"
	POSTGRES_SERVER_HOST = "POSTGRES_SERVER_HOST"
	ENVIRONEMT           = "ENV"
)

type ConfigType map[string]string

var Config = ConfigType{
	POSTGRES_USER:        "",
	POSTGRES_PASSWORD:    "",
	POSTGRES_DB:          "",
	CLIENT_URL:           "",
	SERVER_PORT:          "",
	JWT_KEY:              "",
	RUN_MIGRATION:        "",
	POSTGRES_SERVER_HOST: "localhost",
}

func InitConfig() {
	environment, exists := os.LookupEnv(ENVIRONEMT)
	var envFilePath string
	if exists && environment == "test" {
		envFilePath, _ = filepath.Abs("../.env.test")
	} else {
		envFilePath, _ = filepath.Abs("../.env")
	}
	if err := godotenv.Load(envFilePath); err != nil {
		log.WithField("reason", err.Error()).Fatal("No .env file found")
	}

	required := map[string]bool{
		POSTGRES_USER:     true,
		POSTGRES_PASSWORD: true,
		POSTGRES_DB:       true,
		CLIENT_URL:        true,
		SERVER_PORT:       true,
		RUN_MIGRATION:     true,
	}

	for key := range Config {
		envVal, exists := os.LookupEnv(key)
		if !exists {
			if required[key] {
				log.Fatal(key + " not found in env")
			}
			continue
		}
		if _, ok := Config[key]; ok {
			Config[key] = envVal
		}
	}

	log.Info("All config & secrets set")
}
