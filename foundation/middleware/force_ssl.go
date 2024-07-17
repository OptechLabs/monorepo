package middleware

import (
	"github.com/OptechLabs/monorepo/foundation"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// ForceSSL returns a middleware that redirects to https on production
func ForceSSL(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sec := secure.New(secure.Options{
			SSLRedirect:     env == foundation.Production,
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		})
		if err := sec.Process(c.Writer, c.Request); err != nil {
			// if there was an error, do not continue.
			c.Abort()
			return
		}
		if status := c.Writer.Status(); status > 300 && status < 399 {
			// avoid header rewrite if response is a redirection.
			c.Abort()
		}
	}
}
