package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseSubdomain(environment string, localSubdomainDefault string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subdomain, err := getSubdomain(c.Request.Host)
		if err != nil {
			if environment == "development" || environment == "test" {
				subdomain = localSubdomainDefault
			} else {
				if !strings.Contains(c.Request.URL.Path, "/select_subdomain") {
					c.HTML(http.StatusNotFound, "errors/404", nil)
				}
			}
		}
		c.Set("subdomain", subdomain)
		c.Next()
	}
}

func getSubdomain(host string) (string, error) {
	parts := strings.Split(host, ".")
	if len(parts) != 3 {
		return "", errors.New("no subdomain provided")
	}
	return parts[0], nil
}
