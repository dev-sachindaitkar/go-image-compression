package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateFileSize(maxFileSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			contentLengthStr := c.GetHeader("Content-Length")

			if contentLengthStr != "" {
				contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid content length header format"})
					return
				}
				if contentLength > maxFileSize {
					c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "Payload too large. Maximum supported size is 10MB.",
					})
					return // Halts further execution pipeline stages
				}
			}
		}
		c.Next()
	}
}
