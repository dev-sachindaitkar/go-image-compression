// Fallback guarantees that if the .env.local file is unreadable, it won't crash your server
export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:9090";

export const API_ROUTES = {
  UPLOAD: `${API_BASE_URL}/api/upload`,
  PROGRESS: `${API_BASE_URL}/api/progress`,
  DOWNLOAD: (id: string) => `${API_BASE_URL}/api/download/${id}`,
};

export const MAX_FILE_SIZE_MB = 10;
