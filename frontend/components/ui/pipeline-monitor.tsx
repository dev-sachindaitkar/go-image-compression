"use client";

import React from "react";
import { FileProgress } from "@/types";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { API_ROUTES } from "@/constants/constants";

interface PipelineMonitorProps {
  activeJobs: FileProgress[];
}

export function PipelineMonitor({ activeJobs }: PipelineMonitorProps) {
  const formatKB = (bytes?: number) =>
    bytes ? `${(bytes / 1024).toFixed(1)} KB` : "0 KB";

  return (
    /* Matches your other components' base card style, but allowed to fill the width grid space */
    <div className="bg-white/80 backdrop-blur-md rounded-xl shadow-soft p-6 w-full h-full flex flex-col overflow-hidden border border-slate-100/60">
      {/* Structural Uppercase Header matching CompressionSettings/BatchInput */}
      <h2 className="text-sm font-semibold text-slate-700 uppercase mb-6">
        Live Pipeline Activity Monitoring
      </h2>

      <div className="flex-1 overflow-y-auto pr-1 min-h-0 custom-scrollbar">
        {activeJobs.length === 0 ? (
          /* Clean, transparent empty state that respects the canvas */
          // <div className="flex flex-col items-center justify-center h-64 border border-dashed rounded-xl border-slate-200 text-slate-400">
          <div className="bg-white/80 backdrop-blur-md rounded-xl shadow-soft p-6 w-full h-full flex flex-col overflow-hidden border border-slate-100">
            <p className="text-sm font-medium font-mono">
              No active payloads running in pipeline queue.
            </p>
            <p className="text-xs mt-1 font-mono text-slate-400/70">
              Upload files to watch concurrent worker allocations.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {activeJobs.map((job) => {
              const savings =
                job.originalSize && job.compressedSize
                  ? Math.round(
                      ((job.originalSize - job.compressedSize) /
                        job.originalSize) *
                        100,
                    )
                  : 0;

              const isDone = job.status === "done";

              return (
                <div
                  key={job.jobId}
                  className="p-5 rounded-xl border border-slate-100 bg-white/60 shadow-xs transition-all duration-300"
                >
                  {/* File Metadata Row */}
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-sm font-bold tracking-wide text-slate-700 truncate max-w-[180px] sm:max-w-xs font-mono">
                      {job.filename}
                    </span>
                    <span
                      className={`text-[10px] uppercase font-mono tracking-wider px-2.5 py-0.5 font-bold rounded ${
                        isDone
                          ? "bg-sky-500/10 text-sky-600"
                          : job.status === "processing"
                            ? "bg-sky-500/10 text-sky-600 animate-pulse"
                            : job.status === "failed"
                              ? "bg-rose-500/10 text-rose-500"
                              : "bg-slate-100 text-slate-500"
                      }`}
                    >
                      {isDone ? "COMPLETED" : job.status}
                    </span>
                  </div>

                  {isDone ? (
                    /* Clean, flush sizing breakdown layout */
                    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between pt-1 gap-3">
                      <div className="text-xs space-y-1 text-slate-500 font-mono">
                        <div>
                          Original Size:{" "}
                          <span className="text-slate-700 font-medium">
                            {formatKB(job.originalSize)}
                          </span>
                        </div>
                        <div>
                          Optimized Size:{" "}
                          <span className="text-sky-600 font-bold">
                            {formatKB(job.compressedSize)}
                          </span>
                        </div>
                        {savings > 0 && (
                          <div className="text-emerald-600 font-extrabold text-[11px] bg-emerald-500/10 px-2 py-0.5 rounded inline-block mt-1">
                            -{savings}%
                          </div>
                        )}
                      </div>

                      {/* Premium borderless sky-blue solid block action trigger */}
                      <Button
                        asChild
                        size="sm"
                        className="w-full sm:w-auto bg-sky-500 hover:bg-sky-600 text-white font-bold tracking-wide transition-all rounded-md border-none shadow-xs cursor-pointer px-4"
                      >
                        <a href={API_ROUTES.DOWNLOAD(job.jobId)} download>
                          Download
                        </a>
                      </Button>
                    </div>
                  ) : (
                    /* Active Processing States */
                    <div className="space-y-2.5 w-full">
                      <Progress
                        // Explicitly bind to absolute numerical zero if queued
                        value={job.status === "queued" ? 0 : job.percentage}
                        className="h-3 bg-slate-100/80 rounded-full overflow-hidden transition-all duration-300"
                      />
                      <div className="flex justify-between items-center text-[11px] text-slate-400 font-mono">
                        <span>ID: {job.jobId.substring(0, 8)}</span>
                        <span className="font-bold transition-colors text-xs">
                          {job.status === "queued" ? "0" : job.percentage}%
                        </span>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
