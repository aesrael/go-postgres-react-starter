package api

import (
	"goapp/packages/config"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthorizeSession(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.SendStatus(http.StatusUnauthorized)
	}

	claims, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config[config.JWT_KEY]), nil
	})
	c.Locals("user", claims)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}
	return c.Next()
}
