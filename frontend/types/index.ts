export interface FileProgress {
  jobId: string;
  filename: string;
  percentage: number;
  status: "queued" | "processing" | "done" | "failed";
  originalSize?: number;
  compressedSize?: number;
}

export interface UploadResponse {
  message: string;
  batch_ids: string[];
}
