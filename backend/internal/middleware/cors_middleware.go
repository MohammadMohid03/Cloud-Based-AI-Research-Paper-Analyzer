package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) headers.
// This is required for the frontend (running on a different port/domain)
// to communicate with the backend API.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the request origin is in the allowed list
		allowed := false
		for _, o := range allowedOrigins {
			if strings.TrimSpace(o) == origin || o == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			// Set CORS headers for allowed origins
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours
		}

		// Handle preflight OPTIONS requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
