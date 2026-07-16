package handlers

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	DownloadCache map[string]CachedAssets
	CacheMu       sync.RWMutex
}

func NewDownloadHandler(cleanupInterval, ttl time.Duration) *DownloadHandler {
	h := &DownloadHandler{
		DownloadCache: make(map[string]CachedAssets),
	}
	go h.startCacheEvictionLoop(cleanupInterval, ttl)
	return h
}

func (h *DownloadHandler) startCacheEvictionLoop(interval, ttl time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		h.CacheMu.Lock()
		now := time.Now()
		evicted := 0
		for id, asset := range h.DownloadCache {
			if now.Sub(asset.CreatedAt) > ttl {
				delete(h.DownloadCache, id)
				evicted++
			}
		}
		if evicted > 0 {
			slog.Warn("[Cache Eviction] Cleaned expired files from cache", "count", evicted)
		}
		h.CacheMu.Unlock()
	}
}

func (h *DownloadHandler) DownloadFile(c *gin.Context) {
	jobId := c.Param("id")
	if jobId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID parameter is required"})
		return
	}
	h.CacheMu.RLock()
	asset, exists := h.DownloadCache[jobId]
	cacheSize := len(h.DownloadCache)
	h.CacheMu.RUnlock()

	if !exists {
		slog.Error(
			"[Download Failed] Target JobId not found in memory cache map",
			"requested_id", jobId,
			"total_cached_items", cacheSize,
		)
		c.JSON(http.StatusNotFound, gin.H{"error ": "File not found"})
		return
	}
	slog.Info("[Download Success] Serving optimized binary stream to client", "job_id", jobId, "filename", asset.Result.FileName)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=optimized_"+asset.Result.FileName)
	c.Data(http.StatusOK, "image/jpeg", asset.Result.CompressedData)
}
