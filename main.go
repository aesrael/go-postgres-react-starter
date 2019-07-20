package main

import (
	"fmt"
	controller "go-postgre-jwt-boilerplate/controller"

	"go-postgre-jwt-boilerplate/errors"

	"go-postgre-jwt-boilerplate/db"

	_ "github.com/dgrijalva/jwt-go"
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
	// router.Use(middlewares.Connect)
	// router.Use(middlewares.ErrorHandler)

	// Statics
	// router.Static("/public", "./public")
	// Ping test
	router.GET("/ping", controller.Pong)
	return router
}

func main() {
	rows, err := db.DB.Query("SELECT * FROM users")
	defer rows.Close()
	errors.HandleErr(err)
	fmt.Println(rows)
	// for rows.Next() {
	// 	var title string
	// 	if err := rows.Scan(&title); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(title)
	// }
	// r := setupRouter()
	// Listen and Serve in 0.0.0.0:8080
	// r.Run(":8081")
}
