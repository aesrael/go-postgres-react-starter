package api

import (
	"goapp/packages/config"
	"goapp/packages/db"

	"github.com/apex/log"
)

func StartServer() {
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Db connection error occurred")
	}
	defer conn.Close()

	runMigration := config.Config["RUN_MIGRATION"] == "true"
	if runMigration && conn != nil {
		if err := db.Migrate(conn, config.Config["DB_NAME"]); err != nil {
			log.WithField("reason", err.Error()).Fatal("db migration failed")
		}
	}

	server := httpServer(conn)
	server.Listen(config.Config["SERVER_PORT"])
}
