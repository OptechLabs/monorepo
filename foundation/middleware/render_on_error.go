package middleware

import (
	"github.com/gin-gonic/gin"
)

func RenderHTMLOnError(errorTemplateMap map[int]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Writer.Status()
		if _, found := errorTemplateMap[status]; found {
			c.HTML(status, errorTemplateMap[status], nil)
		}
		c.Next()
	}
}
