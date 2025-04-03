package middleware

import "github.com/gin-gonic/gin"

func APIKEYAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for the API key in the request header
		if c.GetHeader("x-api-key") != apiKey {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
