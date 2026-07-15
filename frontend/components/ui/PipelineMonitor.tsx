"use client";

import React from "react";
import { FileProgress } from "@/types";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { API_ROUTES } from "@/constants/constants";

interface PipelineMonitorProps {
  activeJobs: FileProgress[];
}

export function PipelineMonitor({ activeJobs }: PipelineMonitorProps) {
  const formatKB = (bytes?: number) =>
    bytes ? `${(bytes / 1024).toFixed(1)} KB` : "0 KB";

  return (
    <Card className="min-h-[380px]">
      <CardHeader>
        <CardTitle className="text-lg font-semibold text-foreground">
          Live Pipeline Activity Monitoring
        </CardTitle>
      </CardHeader>
      <CardContent>
        {activeJobs.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-56 border border-dashed rounded-lg bg-muted/10 text-muted-foreground">
            <p className="text-sm">
              No active payloads running in pipeline queue.
            </p>
            <p className="text-xs mt-1">
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

              return (
                <div
                  key={job.jobId}
                  className="p-4 rounded-xl border bg-card text-card-foreground shadow-sm flex flex-col space-y-3"
                >
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium font-mono truncate max-w-[250px] sm:max-w-xs">
                      {job.filename}
                    </span>
                    <span
                      className={`text-xs px-2.5 py-0.5 font-semibold rounded-full capitalize border ${
                        job.status === "done"
                          ? "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"
                          : job.status === "processing"
                            ? "bg-amber-500/10 text-amber-500 border-amber-500/20"
                            : job.status === "failed"
                              ? "bg-destructive/10 text-destructive border-destructive/20"
                              : "bg-muted text-muted-foreground border-transparent"
                      }`}
                    >
                      {job.status}
                    </span>
                  </div>

                  {job.status === "done" ? (
                    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between bg-muted/40 border p-3 rounded-lg gap-3">
                      <div className="text-xs space-y-0.5 text-muted-foreground font-mono">
                        <div>
                          Original:{" "}
                          <span className="font-semibold text-foreground">
                            {formatKB(job.originalSize)}
                          </span>
                        </div>
                        <div>
                          Optimized:{" "}
                          <span className="font-semibold text-emerald-500">
                            {formatKB(job.compressedSize)}
                          </span>
                        </div>
                        {savings > 0 && (
                          <div className="text-emerald-500 font-bold">
                            Saved {savings}% Space
                          </div>
                        )}
                      </div>
                      <Button
                        asChild
                        size="sm"
                        className="w-full sm:w-auto bg-emerald-600 hover:bg-emerald-700 text-white font-bold border-transparent"
                      >
                        <a href={API_ROUTES.DOWNLOAD(job.jobId)} download>
                          Download Image
                        </a>
                      </Button>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      <Progress value={job.percentage} className="h-2" />
                      <div className="flex justify-between items-center text-xs text-muted-foreground font-mono">
                        <span>ID: {job.jobId.substring(0, 8)}...</span>
                        <span>{job.percentage}%</span>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
