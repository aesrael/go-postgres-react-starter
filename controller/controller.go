package controller

import (
	"fmt"
	"go-postgre-jwt-boilerplate/db"
	"go-postgre-jwt-boilerplate/errors"
	"go-postgre-jwt-boilerplate/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Pong tests that api is working and returning json
func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}

type Claims struct {
	db.User
	jwt.StandardClaims
}

//Create new user
func Create(c *gin.Context) {
	var user db.Register
	c.Bind(&user)
	exists := checkUserExists(user)
	if exists == true {
		errors.ValidationErrors = append(errors.ValidationErrors, "email exists already exists")
	}
	utils.ValidateUser(user, errors.ValidationErrors)

	if len(errors.ValidationErrors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors.ValidationErrors})
		return
	}
	HashPassword(&user)
	_, err := db.DB.Query(db.CreateUserQuery, user.Name, user.Password, user.Email)
	errors.HandleErr(err)
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "user created succesfully"})
}

func Login(c *gin.Context) {
	var jwtKey = []byte("secret")

	var user db.Login
	c.Bind(&user)

	rows, err := db.DB.Query(db.LoginQuery, user.Email)
	errors.HandleErr(err)
	for rows.Next() {

		var id int
		var name, email, password, createdAt, updatedAt string

		err := rows.Scan(&id, &name, &password, &email, &createdAt, &updatedAt)
		errors.HandleErr(err)

		match := CheckPasswordHash(user.Password, password)
		if !match {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "incorrect credentials"})
			return
		}

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Hour)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{

			User: db.User{
				Name: name, Email: email, CreatedAt: createdAt, UpdatedAt: updatedAt,
			},
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		errors.HandleErr(err)
		// c.SetCookie("token", tokenString, expirationTime, "", "*", true, false)
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		fmt.Println(tokenString, token, "get here")
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": "logged in succesfully", "user": claims.User, "token": tokenString})
	}
}

func HashPassword(user *db.Register) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	errors.HandleErr(err)
	user.Password = string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func checkUserExists(user db.Register) bool {
	rows, err := db.DB.Query(db.CheckUserExists, user.Email)
	errors.HandleErr(err)
	if !rows.Next() {
		return false
	}
	return true
}
