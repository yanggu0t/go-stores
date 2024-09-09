package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/config"
)

func Response(c *gin.Context, httpStatus int, status string, msgID string, data interface{}) {
	lang, _ := c.Get("language")
	language, ok := lang.(string)
	if !ok {
		language = "en"
	}

	translatedMsg := config.Translate(msgID, language, nil)

	c.JSON(httpStatus, gin.H{
		"status": status,
		"msg":    translatedMsg,
		"data":   data,
	})
}
