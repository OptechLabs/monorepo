package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getSubdomain(t *testing.T) {
	t.Parallel()
	// Arrange
	host := "wu.otosapp.com"
	// Act
	result, err := getSubdomain(host)
	assert.NoError(t, err)
	// Assert
	assert.Equal(t, result, "wu")
}
