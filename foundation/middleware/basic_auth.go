package middleware

import (
	"net/http"

	"github.com/OptechLabs/monorepo/foundation"
	"github.com/gin-gonic/gin"
)

// BasicAuth returns a middleware that allows you to use basic authentication...not for real production
func BasicAuth(environment string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if environment != "development" {
			logger := foundation.LoggerFrom(c)
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				logger.Warn("no basic auth present")
				c.Header("WWW-Authenticate", `Basic realm="Give username and password"`)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No basic auth present"})
				return
			}

			if !isAuthorised(username, password) {
				logger.Warn("incorrect basic auth present")
				c.Header("WWW-Authenticate", `Basic realm="Give username and password"`)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No basic auth present"})
				return
			}
			logger.Info("basic auth successful")
		}
		c.Next()
	}
}

func isAuthorised(username, password string) bool {
	users := map[string]string{
		"optech":  "ProtectYoNeck",
		"barrett": "Since1941",
		"tpj":     "UncrushEm",
	}

	pass, ok := users[username]
	if !ok {
		return false
	}

	return password == pass
}
