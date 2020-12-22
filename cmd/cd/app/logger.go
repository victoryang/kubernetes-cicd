package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger formatter
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.RequestURI

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		username, _ := c.Get("username")
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		info, _ := c.Get("info")

		comment := c.Errors.String()
		if strings.HasSuffix(c.Request.URL.Path, ".json") {
			fmt.Fprintf(os.Stdout, "%v|%d|%v|%v|%s|%s|%s|%v|%v\n",
				end.Format("2006-01-02 15:04:05"),
				statusCode,
				latency,
				username,
				clientIP,
				method,
				path,
				info,
				strings.TrimSpace(comment),
			)
		} else {
			fmt.Fprintf(os.Stdout, "%v|%d|%v|%s|%s|%s|%v\n",
				end.Format("2006-01-02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				strings.TrimSpace(comment),
			)
		}

	}
}
