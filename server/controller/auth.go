package controller

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks that token is valid, see https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
func AuthMiddleware(c *gin.Context, jwtKey []byte) (jwt.MapClaims, bool) {
	//obtain session token from the requests cookies
	ck, err := c.Request.Cookie("token")
	fmt.Println(ck, "coookie")
	if err != nil {
		fmt.Print(err)
		return nil, false
	}

	// Get the JWT string from the cookie
	tokenString := ck.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	}
	return nil, false
}
