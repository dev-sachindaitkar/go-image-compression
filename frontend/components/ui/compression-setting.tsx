import { startTransition } from "react";
import { Slider } from "./slider";

interface CompressionSettingProps {
  compressionQuality: number[];
  onCompressionSettingChange: (val: number[]) => void;
}

export default function CompressionSettings({
  compressionQuality,
  onCompressionSettingChange,
}: CompressionSettingProps) {
  return (
    <div className="bg-white/80 backdrop-blur-md rounded-xl shadow-soft p-6 w-full">
      <h2 className="text-sm font-semibold text-slate-700 uppercase mb-6">
        Compression Settings
      </h2>
      <div className="flex items-center justify-between mb-4">
        <span className="text-slate-600 font-medium text-sm">
          Target Quality
        </span>
        <span className="text-sky-600 font-bold bg-sky-50 px-2 py-0.5 rounded border border-sky-100 text-sm font-mono">
          {compressionQuality[0]}%
        </span>
      </div>

      <Slider
        max={100}
        min={1}
        step={1}
        value={compressionQuality}
        onValueChange={(value) => {
          startTransition(() => {
            if (Array.isArray(value)) {
              onCompressionSettingChange([...value]);
            } else if (typeof value === "number") {
              onCompressionSettingChange([value]);
            }
          });
        }}
        className="w-full"
      />
    </div>
  );
}
