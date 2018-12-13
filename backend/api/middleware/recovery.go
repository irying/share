package middleware

import (
	"backend/utils"
	"io"
	"net/http/httputil"

	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"
)

// Recovery recover form panic
func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return RecoveryWithWriter(f, gin.DefaultErrorWriter)
}

// RecoveryWithWriter log the panic
func RecoveryWithWriter(f func(c *gin.Context, err interface{}), out io.Writer) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 3)
				reset := string([]byte{27, 91, 48, 109})
				utils.LogPanic("[Nice Recovery] panic recovered:\n\n%s%s\n\n%s%s", httprequest, goErr.Error(), goErr.Stack(), reset)

				f(c, err)
			}
		}()
		c.Next() // execute all the handlers
	}
}
