package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const MaxFileSize = 10 * 1024 * 1024

func ValidateFileSize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// only validating POST requests
		if c.Request.Method == http.MethodPost {
			contentLengthStr := c.GetHeader("Content-Length")

			if contentLengthStr != "" {
				contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid content length header format"})
					return
				}
				if MaxFileSize < contentLength {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "Payload too large. Maximum support size is 10MB / image ",
					})
				}
			}

		}
		c.Next()
	}
}
