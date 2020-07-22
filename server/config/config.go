package config

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	DB_USER     = "user"
	DB_PASSWORD = "password"
)
const (
	DB_NAME     = "goauth"
	CLIENT_URL  = "http://localhost:5431"
)

func Init(){
	if err := godotenv.Load("config_var"); err != nil{
		panic("Could not load environment configuration variables")
	}
	DB_USER,_ = os.LookupEnv("DB_USER")
	DB_PASSWORD,_ = os.LookupEnv("DB_PASSWORD")
}
