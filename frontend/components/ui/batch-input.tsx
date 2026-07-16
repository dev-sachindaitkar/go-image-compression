import { ChangeEvent } from "react";
import { Button } from "./button";

interface BatchInputProps {
  onFileUpload: (e: ChangeEvent<HTMLInputElement>) => void;
  isUploading: boolean;
}

export default function BatchInput({
  onFileUpload,
  isUploading,
}: BatchInputProps) {
  return (
    <div className="bg-white/80 backdrop-blur-md rounded-xl shadow-soft p-6 w-full ">
      <div className="flex flex-col items-center">
        <div className="bg-primary/10 rounded-full p-3 mb-3">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6 text-primary"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1M12 12V4m0 0l-3 3m3-3l3 3"
            />
          </svg>
        </div>
        <h2 className="text-sm font-semibold text-slate-700 uppercase mb-1">
          Batch Input
        </h2>
        <p className="text-xs text-slate-500 mb-4">
          Accepts JPEG / PNG images up to 10MB
        </p>
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
            size="sm"
            disabled={isUploading}
            className="cursor-pointer border-border/80 text-foreground hover:bg-primary hover:text-primary-foreground font-semibold tracking-wide transition-all"
          >
            <span>Select Files</span>
          </Button>
        </label>
      </div>
    </div>
  );
}
