package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                        string
	WorkerCount                 int
	MaxFileSizeBytes            int64
	CacheCleanupIntervalMinutes time.Duration
	CacheTTLMinutes             time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Port:                        getEnv("PORT", "8080"),
		WorkerCount:                 getEnvAsInt("OS_WORKER_COUNT", 0),
		MaxFileSizeBytes:            getEnvAsInt64("MAX_FILE_SIZE_BYTES", 10*1024*1024),
		CacheCleanupIntervalMinutes: time.Duration(getEnvAsInt("CACHE_CLEANUP_INTERVAL_MINUTES", 5)) * time.Minute,
		CacheTTLMinutes:             time.Duration(getEnvAsInt("CACHE_TTL_MINUTES", 15)) * time.Minute,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}
	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			return val
		}
	}
	return fallback
}
