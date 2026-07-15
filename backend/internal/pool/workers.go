package pool

import (
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/dev-sachindaitkar/go-image-compression/internal/compression"
)

var OnJobComplete func(Result)

type Job struct {
	Id       string
	FileName string
	Data     []byte
	Format   string
	Quality  int
}

type ProgressUpdate struct {
	JobId          string `json:"job_id"`
	FileName       string `json:"filename"`
	Percentage     int    `json:"percentage"`
	Status         string `json:"status"`                    // "queued", "processing", "done", "failed"
	OriginalSize   int64  `json:"original_size,omitempty"`   // ADD THIS
	CompressedSize int64  `json:"compressed_size,omitempty"` // ADD THIS
}

type Result struct {
	JobId          string
	FileName       string
	CompressedData []byte
	OriginalSize   int64
	CompressedSize int64
	Err            error
}

type Dispatcher struct {
	MaxWorkers int
	JobQueue   chan Job
	Progress   chan ProgressUpdate
	Results    chan Result
	wg         sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		MaxWorkers: maxWorkers,
		// Bounded channels to prevent memory ballooning under sudden heavy load
		JobQueue: make(chan Job, 100),
		Progress: make(chan ProgressUpdate, 200),
		Results:  make(chan Result, 100),
	}
}

func (d *Dispatcher) Start() {
	// 1. Determine Concurrency level : check env overrid, otherwise fallback to default
	workerCount := d.MaxWorkers
	if envWrokers := os.Getenv("OS_WORKER_COUNT"); envWrokers != "" {
		if val, err := strconv.Atoi(envWrokers); err == nil && val > 0 {
			workerCount = val
		}
	}

	if workerCount <= 0 {
		workerCount = runtime.NumCPU() // DEfaut to system capabilities
	}

	slog.Info("[Dispatcher]: ", "Starting worker pool with %d concurrent workers....", workerCount)
	// 2 . Lets launch the background worker routines.
	for i := range workerCount {
		d.wg.Add(1)
		go d.worker(i + 1)
	}
}

func (d *Dispatcher) worker(id int) {
	defer d.wg.Done()
	slog.Info("[Worker", strconv.Itoa(id), "], Standing by for image tasks")

	// rage over a channel automatically blocks and waits until a new job arrives.
	// cleanly exists the when the channel is closed.
	for job := range d.JobQueue {
		log.Printf("[Worker %d] Processing job %s (%s)", id, job.Id, job.FileName)

		// Broadcasting processing status immediately
		d.Progress <- ProgressUpdate{
			JobId:        job.Id,
			FileName:     job.FileName,
			Percentage:   25,
			Status:       "processing",
			OriginalSize: int64(len(job.Data)),
		}

		config := compression.CompressConfig{
			Quality: job.Quality,
		}

		compressedBytes, err := compression.CompressImage(job.Data, config)
		if err != nil {
			log.Printf("[Worker %d] Critical error processing job %s: %v", id, job.Id, err)

			// Broadcast failure state to the user interface
			d.Progress <- ProgressUpdate{
				JobId:      job.Id,
				FileName:   job.FileName,
				Percentage: 100,
				Status:     "failed",
			}

			// Report error up the pipeline to clean up the request context
			d.Results <- Result{
				JobId:    job.Id,
				FileName: job.FileName,
				Err:      err,
			}
			continue
		}

		d.Progress <- ProgressUpdate{
			JobId:          job.Id,
			FileName:       job.FileName,
			Percentage:     50,
			Status:         "processing",
			OriginalSize:   int64(len(job.Data)),
			CompressedSize: int64(len(compressedBytes)),
		}
		res := Result{
			JobId:          job.Id,
			FileName:       job.FileName,
			OriginalSize:   int64(len(job.Data)),
			CompressedSize: int64(len(compressedBytes)),
			Err:            nil,
			CompressedData: compressedBytes,
		}

		time.Sleep(50 * time.Millisecond)

		if OnJobComplete != nil {
			OnJobComplete(res)
		}

		d.Progress <- ProgressUpdate{
			JobId:          job.Id,
			FileName:       job.FileName,
			Percentage:     100, // constant for now
			Status:         "done",
			OriginalSize:   int64(len(job.Data)),
			CompressedSize: int64(len(compressedBytes)),
		}

		select {
		case d.Results <- res:
		default:
		}

		log.Printf("[Worker %d] Completed job %s successfully", id, job.Id)
	}
	log.Printf("[Worker %d] Shutting down cleanly.", id)
}

// Shutdown ensures all workers finish their active tasks before exiting
func (d *Dispatcher) Shutdown() {
	log.Println("[Dispatcher] Initiating graceful shutdown sequence...")
	close(d.JobQueue)
	// Wait for all active worker loops to finish processing their current image and exit
	d.wg.Wait()
	close(d.Progress)
	close(d.Results)
	log.Println("[Dispatcher] All workers stopped. Exiting engine..")
}
