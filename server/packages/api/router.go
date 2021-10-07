package api

import (
	"database/sql"
	"goapp/packages/api/controller"
	"goapp/packages/api/middleware"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config["CLIENT_URL"],
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		AllowMethods:     "POST, OPTIONS, GET, PUT",
	}))

	api := app.Group("/api")

	// public
	api.Get("/ping", controller.Pong)

	api.Post("/login", WithDB(controller.Login, db))
	api.Post("/register", WithDB(controller.CreateUser, db))

	api.Post("/createReset", WithDB(controller.InitiatePasswordReset, db))
	api.Post("/resetPassword", WithDB(controller.ResetPassword, db))

	// authed routes
	api.Get("/session", middleware.AuthorizeSession, WithDB(controller.Session, db))

	return app
}
