package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log stack trace and error
				fmt.Printf("Panic recovered: %v\n", err)
				debug.PrintStack()

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error ",
				})
			}
		}()

		c.Next()
	}
}
