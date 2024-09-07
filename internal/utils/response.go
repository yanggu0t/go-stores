package utils

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpStatus int, status string, msg string, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"status": status,
		"msg":    msg,
		"data":   data,
	})
}