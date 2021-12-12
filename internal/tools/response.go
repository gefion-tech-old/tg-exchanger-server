package tools

import (
	"github.com/gin-gonic/gin"
)

func ServErr(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
	c.Abort()
}
