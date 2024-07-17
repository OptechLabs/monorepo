package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OptechLabs/monorepo/foundation/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_AuthenticationMW(t *testing.T) {
	router := gin.New()
	router.Use(middleware.BasicJWT("token"))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "200 OK")
	})

	tests := []struct {
		name       string
		givenToken string
		wantCode   int
	}{
		{name: "bad token", givenToken: "Bearer bad token", wantCode: http.StatusUnauthorized},
		{name: "no bearer", givenToken: "bad token", wantCode: http.StatusUnauthorized},
		{name: "no token", givenToken: "", wantCode: http.StatusUnauthorized},
		{name: "happy path", givenToken: "Bearer token", wantCode: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", rr.Body)
			if tc.givenToken != "" {
				req.Header.Add("Authorization", tc.givenToken)
			}
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantCode, rr.Code)
		})
	}
}
