"use client";

import React, { useState, ChangeEvent, startTransition } from "react";
import { useProgressStream } from "@/hooks/useProgressStream";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Slider } from "@/components/ui/slider";
import { API_ROUTES } from "@/constants/constants";
import { UploadDropzone } from "@/components/ui/UploadDropzone";
import { PipelineMonitor } from "@/components/ui/PipelineMonitor";
import dynamic from "next/dynamic";
// import { UploadDropzone } from "@/components/UploadDropzone";
// import { PipelineMonitor } from "@/components/PipelineMonitor";
// import { ThemeToggle } from "@/components/ui/ThemeToggle";
// import { API_ROUTES } from "@/config/constants";

const ThemeToggle = dynamic(() => import("@/components/ui/theme-toggle"), {
  ssr: false,
  loading: () => (
    <div className="h-10 w-10 rounded-xl border bg-background animate-pulse" />
  ),
});

export default function Page() {
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

      if (!response.ok) throw new Error("Upload transaction failed");
    } catch (err) {
      console.error("Network dispatch error:", err);
      alert("Failed to transmit files to the backend server matrix.");
    } finally {
      setIsUploading(false);
      e.target.value = "";
    }
  };

  return (
    <main className="min-h-screen bg-background text-foreground px-6 py-12 transition-colors duration-200">
      <div className="mx-auto max-w-5xl">
        <header className="mb-12 border-b pb-6 flex items-center justify-between border-border">
          <div>
            <h1 className="text-3xl font-extrabold tracking-tight">
              GoOptiFlow{" "}
              <span className="text-sm font-medium text-primary px-2.5 py-1 bg-primary/10 rounded-md ml-2 border border-primary/20">
                Design System Architecture
              </span>
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              High-performance parallel image optimization powered by Go.
            </p>
          </div>
          <ThemeToggle />
        </header>

        <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
          <div className="lg:col-span-1 space-y-6">
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-md font-semibold text-foreground">
                  Parameters
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-between items-center text-sm font-medium">
                  <span className="text-muted-foreground">
                    Compression Quality
                  </span>
                  <span className="font-bold text-primary font-mono">
                    {compressionQuality[0]}%
                  </span>
                </div>
                <Slider
                  min={1}
                  max={100}
                  step={1}
                  value={compressionQuality}
                  onValueChange={(value) => {
                    // startTransition prevents cascading render errors in React 19
                    startTransition(() => {
                      if (Array.isArray(value)) {
                        setCompressionQuality([...value]);
                      } else if (typeof value === "number") {
                        setCompressionQuality([value]);
                      }
                    });
                  }}
                  className="cursor-pointer"
                />
              </CardContent>
            </Card>

            <UploadDropzone
              onFileUpload={handleFileUpload}
              isUploading={isUploading}
            />
          </div>

          <div className="lg:col-span-2">
            <PipelineMonitor activeJobs={Object.values(progressMap)} />
          </div>
        </div>
      </div>
    </main>
  );
}
