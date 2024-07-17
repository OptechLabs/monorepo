package foundation

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	I18nKey      = "i18n"
	LoggerKey    = "foundationLogger"
	RequestIDKey = "requestID"
	TxKey        = "tx"
)

type BuffaloValidateError interface {
	Error() string
	String() string
	HasAny() bool
}

func AbortWithError(c *gin.Context, code int, err error) *gin.Error {
	hashmap := gin.H{
		"request_id": RequestIDFrom(c),
	}
	if _, ok := err.(BuffaloValidateError); ok {
		hashmap["validate"] = err
	} else {
		hashmap["error"] = err.Error()
	}
	c.AbortWithStatusJSON(code, hashmap)
	return c.Error(err)
}

func LoggerFrom(c *gin.Context) Logger {
	if maybeALogger, exists := c.Get(LoggerKey); exists {
		if logger, ok := maybeALogger.(Logger); ok {
			return logger
		}
	}
	defaultLogger, _ := NewDefaultLogger("")
	return defaultLogger
}

func RequestIDFrom(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

func TxFrom(c *gin.Context) (tx *sqlx.Tx, ok bool) {
	if maybeTx, exists := c.Get(TxKey); exists {
		tx, ok = maybeTx.(*sqlx.Tx)
	}
	return
}

func TxMustFrom(c *gin.Context) *sqlx.Tx {
	if tx, exists := TxFrom(c); exists {
		return tx
	}
	panic(`"` + TxKey + `" does not exist in context`)
}
