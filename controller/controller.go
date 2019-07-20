package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pong, test that api is working and returning json
func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}
