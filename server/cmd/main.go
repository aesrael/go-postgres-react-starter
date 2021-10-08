package main

import (
	"goapp/packages/api"
	"goapp/packages/config"
)

func main() {
	config.InitConfig()
	api.StartServer()
}
