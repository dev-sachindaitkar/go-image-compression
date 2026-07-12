package main

import (
	"log/slog"
	"net/http"

	"github.com/dev-sachindaitkar/go-image-compression/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.ValidateFileSize())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Healthy"})
	})

	r.POST("/api/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Payload passed the check"})
	})

	slog.Info("Starting Image compressor servcer on 8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
