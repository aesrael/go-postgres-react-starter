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

var jwtKey = []byte("secret")

//Create new user
func Create(c *gin.Context) {
	var user db.Register
	c.Bind(&user)
	exists := checkUserExists(user)

	valErr := utils.ValidateUser(user, errors.ValidationErrors)
	if exists == true {
		valErr = append(valErr, "email already exists")
	}
	fmt.Println(valErr)
	if len(valErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "errors": valErr})
		return
	}
	HashPassword(&user)
	_, err := db.DB.Query(db.CreateUserQuery, user.Name, user.Password, user.Email)
	errors.HandleErr(err)
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "user created succesfully"})
}

func Session(c *gin.Context) {
	user, isAuthenticated := authMiddleware(c)
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "user": user})
}

func Login(c *gin.Context) {
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
		expirationTime := time.Now().Add(5 * time.Second)
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

func authMiddleware(c *gin.Context) (jwt.MapClaims, bool) {
	// We can obtain the session token from the requests cookies, which come with every request
	ck, err := c.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "unauthorized"})
			return nil, false
		}
		// For any other type of error, return a bad request status
		errors.HandleErr(err)
	}

	// Get the JWT string from the cookie
	tokenString := ck.Value

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return nil, false
	}
	return nil, true
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
