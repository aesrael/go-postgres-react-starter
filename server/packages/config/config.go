package config

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/joho/godotenv"
)

type ConfigType map[string]string

var Config = ConfigType{
	"DB_USER":       "",
	"DB_PASSWORD":   "",
	"DB_NAME":       "",
	"CLIENT_URL":    "",
	"SERVER_PORT":   "",
	"JWT_KEY":       "",
	"RUN_MIGRATION": "",
}

func InitConfig() {
	envFilePath, _ := filepath.Abs("../.env")
	if err := godotenv.Load(envFilePath); err != nil {
		log.WithField("reason", err.Error()).Fatal("No .env file found")
	}

	required := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"CLIENT_URL",
		"SERVER_PORT",
		"RUN_MIGRATION",
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
