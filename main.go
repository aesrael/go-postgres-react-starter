package main

import (
	_ "database/sql"
	controller "go-postgre-jwt-boilerplate/controller"
	"go-postgre-jwt-boilerplate/middlewares"

	"go-postgre-jwt-boilerplate/db"

	// _ "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func init() {
	db.Connect()
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	// Middlewares
	router.Use(middlewares.ErrorHandler)

	// Statics
	// router.Static("/public", "./public")
	// Ping test
	router.GET("/ping", controller.Pong)
	router.POST("/register", controller.Create)
	router.POST("/login", controller.Login)
	router.GET("/session", controller.Session)
	return router
}

func main() {
	r := setupRouter()
	// Listen and Serve in 0.0.0.0:80801
	r.Run(":8081")
}
