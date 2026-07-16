"use client";

import React, { useState, ChangeEvent } from "react";
import { useProgressStream } from "@/hooks/useProgressStream";
import { API_ROUTES } from "@/constants/constants";

import { PipelineMonitor } from "@/components/ui/pipeline-monitor";
import BatchInput from "@/components/ui/batch-input";
import CompressionSettings from "@/components/ui/compression-setting";

export default function GoOptiFlowDashboard() {
  const { progressMap } = useProgressStream();
  const [compressionQuality, setCompressionQuality] = useState<number[]>([75]);
  const [isUploading, setIsUploading] = useState<boolean>(false);

  const handleFileUpload = async (e: ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = e.target.files;
    if (!selectedFiles || selectedFiles.length === 0) return;

    setIsUploading(true);

    const formData = new FormData();
    formData.append("quality", compressionQuality[0].toString());

    for (let i = 0; i < selectedFiles.length; i++) {
      formData.append("images", selectedFiles[i]);
    }

    try {
      const response = await fetch(API_ROUTES.UPLOAD, {
        method: "POST",
        body: formData,
      });

      if (!response.ok)
        throw new Error("Upload dispatch sequence rejected by server node");
    } catch (err) {
      console.error("Network network transmission fault:", err);
      alert("Failed to transmit batch files to the parallel Go engine.");
    } finally {
      setIsUploading(false);
      e.target.value = "";
    }
  };

  return (
    <main className="h-screen min-h-screen w-full p-6 md:p-12 bg-transparent flex flex-col overflow-hidden">
      <div className="mx-auto max-w-5xl w-full h-full flex flex-col min-h-0">
        <header className="mb-12 border-b pb-6 flex items-center justify-between border-border/60 shrink-0">
          <div>
            <h1 className="text-4xl font-extrabold tracking-tight text-slate-800">
              GoOptiFlow
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              High-performance parallel image optimization powered by Go.
            </p>
          </div>
        </header>

        {/* items-start decouples the height matching completely! 
          The left side and right side will now size perfectly based on their own content alone.
        */}
        <div className="grid grid-cols-1 gap-8 xl:grid-cols-3 items-start flex-1 min-h-0">
          {/* Left Column: Compression Settings & Batch Input stacked naturally */}
          <div className="xl:col-span-1 grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-1 gap-6">
            <CompressionSettings
              compressionQuality={compressionQuality}
              onCompressionSettingChange={(v) => {
                setCompressionQuality(v);
              }}
            />

            <BatchInput
              onFileUpload={handleFileUpload}
              isUploading={isUploading}
            />
          </div>

          {/* Right Column: Live Monitor Panel sitting side-by-side naturally */}
          <div className="lg:col-span-2 h-full flex flex-col min-h-0">
            <PipelineMonitor activeJobs={Object.values(progressMap)} />
          </div>
        </div>
      </div>
    </main>
  );
}
