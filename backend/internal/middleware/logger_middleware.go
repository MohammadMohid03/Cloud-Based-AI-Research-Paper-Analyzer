package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware provides structured request logging compatible with CloudWatch.
// It logs method, path, status code, latency, and client IP for each request.
// This structured format makes it easy to parse logs in cloud environments.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time
		startTime := time.Now()

		// Process the request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)

		// Get response status
		statusCode := c.Writer.Status()

		// Determine log level based on status code
		logPrefix := "INFO"
		if statusCode >= 400 && statusCode < 500 {
			logPrefix = "WARN"
		} else if statusCode >= 500 {
			logPrefix = "ERROR"
		}

		// Structured log output compatible with CloudWatch Logs
		log.Printf("[%s] %s %s | status=%d | duration=%v | ip=%s | user-agent=%s",
			logPrefix,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			duration,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		// Log any errors that occurred during request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("[ERROR] Request error: %v", err.Error())
			}
		}
	}
}
