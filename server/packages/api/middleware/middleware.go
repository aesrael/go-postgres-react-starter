package middleware

import (
	"goapp/packages/config"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthorizeSession(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")

	if tokenStr == "" {
		c.SendStatus(http.StatusUnauthorized)
		return nil
	}

	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config["JWT_KEY"]), nil
	})

	if err != nil {
		c.Status(http.StatusUnauthorized).SendString(err.Error())
		return err
	}
	return c.Next()
}
