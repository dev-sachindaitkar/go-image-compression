"use client";

import { useEffect, useState } from "react";
import { FileProgress } from "@/types";
import { API_ROUTES } from "@/constants/constants";

export function useProgressStream() {
  const [progressMap, setProgressMap] = useState<Record<string, FileProgress>>(
    {},
  );

  useEffect(() => {
    // Establishes a persistent Server-Sent Events link to our Go service matrix
    const eventSource = new EventSource(API_ROUTES.PROGRESS);

    eventSource.addEventListener("progress", (event: MessageEvent) => {
      try {
        const rawData = JSON.parse(event.data);

        const updatedProgress: FileProgress = {
          jobId: rawData.job_id,
          filename: rawData.filename,
          percentage: rawData.percentage,
          status: rawData.status,
          originalSize: rawData.original_size,
          compressedSize: rawData.compressed_size,
        };

        // Appends or updates the specific file token record dynamically in state
        setProgressMap((prev) => ({
          ...prev,
          [updatedProgress.jobId]: updatedProgress,
        }));
      } catch (err) {
        console.error("Error parsing SSE pipeline payload frame:", err);
      }
    });

    eventSource.onerror = (err) => {
      console.error("SSE Stream channel connectivity interruption:", err);
      eventSource.close();
    };

    // Clean up function to close the network connection when the component unmounts
    return () => {
      eventSource.close();
    };
  }, []);

  return { progressMap, setProgressMap };
}
