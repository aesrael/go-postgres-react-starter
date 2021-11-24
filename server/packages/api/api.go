package api

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"goapp/packages/config"
	"goapp/packages/db"
)

var server *fiber.App

func StartServer() {
	conn, err := db.ConnectDB()
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Db connection error occurred")
	}
	defer conn.Close()

	runMigration := config.Config[config.RUN_MIGRATION]
	dbName := config.Config[config.POSTGRES_DB]
	port := config.Config[config.SERVER_PORT]

	if runMigration == "true" && conn != nil {
		if err := db.Migrate(conn, dbName); err != nil {
			log.WithField("reason", err.Error()).Fatal("db migration failed")
		}
	}

	server = httpServer(conn)
	serverErr := server.Listen(port)
	if serverErr != nil {
		log.WithField("reason", serverErr.Error()).Fatal("Server error")
	}
}

func StopServer() {
	if server != nil {
		err := server.Shutdown()
		if err != nil {
			log.WithField("reason", err.Error()).Fatal("Shutdown server error")
		}
	}
}
