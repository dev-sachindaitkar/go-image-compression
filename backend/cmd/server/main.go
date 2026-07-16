package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-sachindaitkar/go-image-compression/internal/api/handlers"
	"github.com/dev-sachindaitkar/go-image-compression/internal/api/router"
	"github.com/dev-sachindaitkar/go-image-compression/internal/config"
	"github.com/dev-sachindaitkar/go-image-compression/internal/logger"
	"github.com/dev-sachindaitkar/go-image-compression/internal/pool"
	"github.com/joho/godotenv"
)

func main() {
	logHandler := logger.NewColorTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))
	slog.Info("[System] Initializing Config conext...")

	if err := godotenv.Load(); err != nil {
		slog.Info("[Warning] .env file not found, using default values")
	}

	cfg := config.LoadConfig()

	// 1. Initialize Concurrency Dispatcher using config parameters
	dispatcher := pool.NewDispatcher(cfg.WorkerCount)
	dispatcher.Start()

	// 2. Init modulared handler
	imageHandler := handlers.NewImageHandler(dispatcher)
	downloadHandler := handlers.NewDownloadHandler(cfg.CacheCleanupIntervalMinutes, cfg.CacheTTLMinutes)

	pool.OnJobComplete = func(res pool.Result) {
		downloadHandler.CacheMu.Lock()
		downloadHandler.DownloadCache[res.JobId] = handlers.CachedAssets{
			Result:    res,
			CreatedAt: time.Now(),
		} // Matches your lowercase 'JobId' struct tag
		slog.Info("[Cache Debug] Worker inserted job %s directly into instance cache. Current Size: ", res.JobId, len(downloadHandler.DownloadCache))
		downloadHandler.CacheMu.Unlock()
	}

	r := router.SetupRouter(cfg, imageHandler, downloadHandler)

	// 4. Implement Production Graceful Shutdown Handling
	// This captures structural operating system interrupt flags (Ctrl+C, termination triggers)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Create channel listening for kernel kill signals
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("[System] Listen and serve crash event: ", "error", err.Error())
		}
	}()
	slog.Info("[System] Network port matrix bound successfully on :", cfg.Port, ".")

	// Set up OS signal interception channels
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("[System] Shutdown broadcast signal received. Stopping application safely...")

	// Create a hard 5-second maximum context safety net window for remaining network jobs to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("[System] Server forced to close prematurely: %v", "error", err.Error())
	}

	// Shut down our internal worker pool loops safely
	dispatcher.Shutdown()
}
