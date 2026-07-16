package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dev-sachindaitkar/go-image-compression/internal/pool"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CachedAssets struct {
	Result    pool.Result
	CreatedAt time.Time
}
type ImageHandler struct {
	Dispatcher *pool.Dispatcher
	StagedJobs map[string]pool.Job
}

func NewImageHandler(d *pool.Dispatcher) *ImageHandler {
	h := &ImageHandler{
		Dispatcher: d,
		StagedJobs: make(map[string]pool.Job),
	}
	return h
}

func (h *ImageHandler) UploadBatch(c *gin.Context) {
	// extract form configurations e.g image quality
	qualityStr := c.DefaultPostForm("quality", "75")
	quality, err := strconv.Atoi(qualityStr)
	if err != nil || quality < 1 || quality > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quality parameter must be an integer between 1 and 100"})
		return
	}

	// parse the mutlipart form file array
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to prase multipart upload payload"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files found under key 'images'"})
		return
	}

	type StagedJobResponse struct {
		JobId        string `json:"jobId"`
		Filename     string `json:"filename"`
		Status       string `json:"status"` // will be "queued"
		Percentage   int    `json:"percentage"`
		OriginalSize int64  `json:"originalSize"`
	}

	var responsePayload []StagedJobResponse

	var trackingIDs []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		fileBytes, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			continue
		}

		// Generate a unique tracking ID for this specific image file asset
		jobId := uuid.New().String()
		trackingIDs = append(trackingIDs, jobId)

		stagedJobs := pool.Job{
			Id:       jobId,
			FileName: fileHeader.Filename,
			Data:     fileBytes,
			Quality:  quality,
		}
		h.StagedJobs[jobId] = stagedJobs
		slog.Info("[Staging Engine] Asset staged successfully. Awaiting trigger event...", "job_id", jobId, "filename", fileHeader.Filename)
		responsePayload = append(responsePayload, StagedJobResponse{
			JobId:        jobId,
			Filename:     fileHeader.Filename,
			Status:       "queued", // Staged status on initialization
			Percentage:   0,
			OriginalSize: int64(len(fileBytes)),
		})

		// construct pool job
		// job := pool.Job{
		// 	Id:       jobId,
		// 	FileName: fileHeader.Filename,
		// 	Data:     fileBytes,
		// 	Quality:  quality,
		// }

		// select {
		// case h.Dispatcher.JobQueue <- job:
		// 	h.Dispatcher.Progress <- pool.ProgressUpdate{
		// 		JobId:      jobId,
		// 		FileName:   fileHeader.Filename,
		// 		Percentage: 0,
		// 		Status:     "queued",
		// 	}

		// default:
		// 	// Backpressure fallback if the channel buffer hits absolute overflow limits
		// 	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Server job queue is currently full. Try again later."})
		// 	return
		// }

	}
	// 4. Respond instantly with 202 Accepted status and tracking metadata
	c.JSON(http.StatusOK, responsePayload)
}

// StreamProgress handles long-lived Server-Sent Events (SSE) streaming connections
func (h *ImageHandler) StreamProgress(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allows smooth cross-origin connections

	// Continuously stream channel items down the network socket pipe until client disconnects
	c.Stream(func(w io.Writer) bool {
		// Read a progress update item out of our global worker progress channel
		if update, ok := <-h.Dispatcher.Progress; ok {
			// Write formatted SSE payload string syntax: data: { ...json... }\n\n
			c.SSEvent("progress", update)
			return true // Keep connection stream context alive
		}
		return false // Stops stream tracking if the progress channel gets closed natively
	})
}

func (h *ImageHandler) TriggerCompression(c *gin.Context) {
	type TriggerPayload struct {
		JobIds []string `json:"jobIds"`
	}

	var payload TriggerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
	}

	if len(payload.JobIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No job tracking IDs provided for execution"})
		return
	}

	triggeredCount := 0

	for _, id := range payload.JobIds {
		job, exists := h.StagedJobs[id]

		if !exists {
			slog.Warn("[Pipeline] Attempted to compress a job that is missing or already processed", "job_id", id)
			continue
		}

		select {
		case h.Dispatcher.JobQueue <- job:
			triggeredCount++
			delete(h.StagedJobs, id)
			slog.Info("[Pipeline] Successfully dispatched staged asset to core workers", "job_id", id)
		default:
			slog.Error("[Pipeline] Queue backpressure overflow during execution routing", "job_id", id)
		}

	}
	c.JSON(http.StatusOK, gin.H{
		"message":         "Compression sequence initialized successfully.",
		"triggered_count": triggeredCount,
	})
}
