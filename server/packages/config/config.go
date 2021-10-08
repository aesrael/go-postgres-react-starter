package config

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/joho/godotenv"
)

type ConfigType map[string]string

var Config = ConfigType{
	"DB_USER":       "aesrael",
	"DB_PASSWORD":   "",
	"DB_NAME":       "goapp",
	"CLIENT_URL":    "",
	"SERVER_PORT":   "",
	"RUN_MIGRATION": "true",
}

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.WithField("reason", err.Error()).Fatal("No .env file found")
	}

	a, err := os.LookupEnv("DB_USER")
	fmt.Print(a, err)

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
