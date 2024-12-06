package handler

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	email := c.MustGet("email").(string)
	c.JSON(200, gin.H{"message": "Welcome " + email})
}
