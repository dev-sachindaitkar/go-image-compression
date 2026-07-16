"use client";

import { useEffect, useState } from "react";
import { FileProgress } from "@/types";
import { API_ROUTES } from "@/constants/constants";

export function useProgressStream() {
  const [progressMap, setProgressMap] = useState<Record<string, FileProgress>>(
    {},
  );

  useEffect(() => {
    // Open a persistent connection to the Go background pipeline
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

        // Dynamically add or merge the progress token frame into our state mapping
        setProgressMap((prev) => ({
          ...prev,
          [updatedProgress.jobId]: updatedProgress,
        }));
      } catch (err) {
        console.error("Error parsing incoming SSE layout frame:", err);
      }
    });

    eventSource.onerror = (err) => {
      console.error(
        "SSE connection interrupted. Disconnecting stream channel:",
        err,
      );
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return { progressMap };
}
