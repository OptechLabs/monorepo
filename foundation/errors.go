package foundation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContextErrors struct {
	errs []*gin.Error
}

func NewContextErrors(errs []*gin.Error) *ContextErrors {
	return &ContextErrors{
		errs: errs,
	}
}

func (c *ContextErrors) Error() string {
	if len(c.errs) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, err := range c.errs {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, err.Err)
		if err.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", err.Meta)
		}
	}
	return buffer.String()
}

func (c *ContextErrors) Is(target error) bool {
	if c == target {
		return true
	}
	if ginErr, ok := target.(*gin.Error); ok {
		target = ginErr.Unwrap()
	}
	for _, err := range c.errs {
		if errors.Is(err.Unwrap(), target) {
			return true
		}
	}
	return false
}
