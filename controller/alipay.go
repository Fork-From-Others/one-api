package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(c *gin.Context) {
	data := make(map[string]interface{})
	data["name"] = "Tom"
	data["age"] = 18
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hello World",
		"data":    data,
	})
	return
}
