package config

import (
	"os"

	"github.com/apex/log"
	"github.com/joho/godotenv"
)

type ConfigType map[string]string

var Config = ConfigType{
	"DB_USER":     "",
	"DB_PASSWORD": "",
	"DB_NAME":     "",
	"CLIENT_URL":  "",
	"SERVER_PORT": "",
}

const ALLOWED_ORIGINS = "https://github.com"

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found")
	}

	required := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"CLIENT_URL",
		"SERVER_PORT",
	}

	for _, env := range required {
		envVal, exists := os.LookupEnv(env)
		if !exists {
			log.Fatal(env + " not found in env")
		}
		if _, ok := Config[env]; ok {
			Config[env] = envVal
		}
	}
	log.Info("All config & secrets set")
}
