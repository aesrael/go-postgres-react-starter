package router

import (
	"go-postgres-jwt-react-starter/controller"
	"go-postgres-jwt-react-starter/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRouter setup routing here
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Middlewares
	router.Use(middlewares.ErrorHandler)

	router.GET("/ping", controller.Pong)
	router.POST("/register", controller.Create)
	router.POST("/login", controller.Login)
	router.GET("/session", controller.Session)
	return router
}
