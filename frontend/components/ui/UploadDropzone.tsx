"use client";

import React, { ChangeEvent } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface UploadDropzoneProps {
  onFileUpload: (e: ChangeEvent<HTMLInputElement>) => void;
  isUploading: boolean;
}

export function UploadDropzone({
  onFileUpload,
  isUploading,
}: UploadDropzoneProps) {
  return (
    <Card className="border-dashed border-2 bg-muted/20 hover:bg-muted/40 transition-colors duration-200">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 text-primary mb-2">
          <svg
            className="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
            />
          </svg>
        </div>
        <CardTitle className="text-md font-semibold text-foreground">
          {isUploading ? "Transmitting Batch..." : "Upload Image Batch"}
        </CardTitle>
        <CardDescription className="text-xs">
          Supports JPEG and PNG formats up to 10MB per file
        </CardDescription>
      </CardHeader>
      <CardContent className="flex justify-center pt-2">
        <label className="cursor-pointer">
          <input
            type="file"
            multiple
            accept="image/jpeg,image/png"
            disabled={isUploading}
            onChange={onFileUpload}
            className="hidden"
          />
          <Button
            asChild
            variant="outline"
            disabled={isUploading}
            className="cursor-pointer"
          >
            <span>Select Files</span>
          </Button>
        </label>
      </CardContent>
    </Card>
  );
}
