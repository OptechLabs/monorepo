package middleware

import (
	"net/http"
	"strings"

	"github.com/OptechLabs/monorepo/foundation"
	"github.com/gin-gonic/gin"
)

// BasicJWT returns a middleware that allows you to create a middleware with the SystemToken
func BasicJWT(systemToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := foundation.LoggerFrom(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("no bearer token in request header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			logger.Warn("not a bearer token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		tokenString := splitToken[1]

		if tokenString != systemToken {
			logger.Warn("bearer token not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.Next()
	}
}
