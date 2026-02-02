package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware provides request logging using the decoupled logger
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// Create middleware-specific logger
	middlewareLogger := logger.With(
		zap.String("component", "api_middleware"),
	)

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log using our structured logger instead of gin's default
		logFields := []zap.Field{
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("protocol", param.Request.Proto),
			zap.Int("status_code", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		}

		if param.ErrorMessage != "" {
			logFields = append(logFields, zap.String("error", param.ErrorMessage))
		}

		// Log at appropriate level based on status code
		if param.StatusCode >= 500 {
			middlewareLogger.Error("HTTP request completed with server error", logFields...)
		} else if param.StatusCode >= 400 {
			middlewareLogger.Warn("HTTP request completed with client error", logFields...)
		} else {
			middlewareLogger.Info("HTTP request completed", logFields...)
		}

		// Return empty string since we're handling logging ourselves
		return ""
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles errors
func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	middlewareLogger := logger.With(
		zap.String("component", "error_middleware"),
	)

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			middlewareLogger.Error("Request completed with errors",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("error", err.Error()),
				zap.Uint64("error_type", uint64(err.Type)))
		}
	}
}

// RateLimitMiddleware provides basic rate limiting (placeholder)
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add rate limiting logic here
		c.Next()
	}
}

// AuthMiddleware provides authentication (placeholder)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add authentication logic here
		c.Next()
	}
}
