package api

import (
	"database/sql"
	"goapp/packages/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func WithDB(fn func(c *fiber.Ctx, db *sql.DB) error, db *sql.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return fn(c, db)
	}
}

func httpServer(db *sql.DB) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(requestid.New())

	api := app.Group("/api")
	api.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config[config.CLIENT_URL],
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin",
		AllowMethods:     "POST, OPTIONS, GET, PUT",
		ExposeHeaders:    "Set-Cookie",
	}))

	// public
	api.Get("/ping", Pong)

	api.Post("/login", WithDB(Login, db))
	api.Post("/register", WithDB(CreateUser, db))
	api.Get("/logout", Logout)

	// authed routes
	api.Get("/session", AuthorizeSession, WithDB(Session, db))

	return app
}
