package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_RenderHTMLOnError(t *testing.T) {
	t.Parallel()
	// Arrange
	router := gin.Default()
	router.Use(RenderHTMLOnError(map[int]string{http.StatusNotFound: "404"}))
	router.LoadHTMLFiles("404.html")

	// Act
	router.GET("/not", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/not_here", nil))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusNotFound)
	assert.Contains(t, recorder.Body.String(), "404 page not found")

}
