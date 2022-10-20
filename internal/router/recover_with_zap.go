package router

import (
	"github.com/helegehe/mini_app/tools/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Warnf("[broken pipe] [err:] %v \n [req:] %v", err, c.FullPath())
					// If the connection is dead, we can't write a status to it.
					// nolint: errcheck
					c.Abort()
					return
				}

				logger.Errorf("[Recovery from panic] [time:] %s \n [err:] %s \n [req:] %v \n [stack:] %s", time.Now().String(), err, string(httpRequest), string(debug.Stack()))
				c.Abort()
				resp := make(map[string]string, 0)
				resp["status"] = "failure"
				resp["message"] = "InternalServerError"
				c.JSON(http.StatusInternalServerError, resp)
			}
		}()
		c.Next()
	}
}
