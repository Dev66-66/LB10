package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		fmt.Printf(
			"time=%s method=%s path=%s status=%d latency=%s client=%s\n",
			start.Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start).String(),
			c.ClientIP(),
		)
	}
}
