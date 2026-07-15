package router

import (
	"time"

	"github.com/dev-sachindaitkar/go-image-compression/internal/api/handlers"
	"github.com/dev-sachindaitkar/go-image-compression/internal/api/middleware"
	"github.com/dev-sachindaitkar/go-image-compression/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, imageHandler *handlers.ImageHandler, downloadHandler *handlers.DownloadHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.ValidateFileSize(cfg.MaxFileSizeBytes))
	r.GET("/heath", func(c *gin.Context) {
		c.JSON(200, gin.H{"Status": "Healthy"})
	})

	api := r.Group("/api")
	{
		api.POST("/upload", imageHandler.UploadBatch)
		api.GET("/progress", imageHandler.StreamProgress)
		api.POST("/download/:id", downloadHandler.DownloadFile)
	}
	return r
}
