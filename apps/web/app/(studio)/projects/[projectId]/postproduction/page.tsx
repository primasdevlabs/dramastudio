"use client";

import { useState } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import { Sliders, CheckCircle2, Cpu } from "lucide-react";
import { api } from "@/lib/api/client";
import { RenderTask } from "@/lib/api/types";

export default function PostproductionPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const [renderResult, setRenderResult] = useState<RenderTask | null>(null);

  const renderMutation = useMutation({
    mutationFn: async () => {
      return api.post<RenderTask>("/v1/postproduction/renders", {
        episode_id: "ep_001",
        format: "mp4",
        resolution: "1080x1920",
      });
    },
    onSuccess: (data) => {
      setRenderResult(data);
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-md">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded bg-studio-accent/10 border border-studio-accent/20 flex items-center justify-center text-studio-accent">
            <Sliders className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Postproduction Assembly</h1>
            <p className="text-xs text-studio-muted">Edit Decision List, multi-track timeline assembly, and 9:16 FFmpeg rendering engine</p>
          </div>
        </div>

        <button
          onClick={() => renderMutation.mutate()}
          disabled={renderMutation.isPending}
          className="inline-flex items-center gap-2 bg-studio-accent hover:bg-studio-accent-dark text-white font-medium text-xs px-4 py-2.5 rounded transition-colors"
        >
          <Cpu className="w-4 h-4" />
          {renderMutation.isPending ? "Rendering master..." : "Render Master (FFmpeg 9:16)"}
        </button>
      </div>

      {renderResult && (
        <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 p-4 rounded text-xs flex items-center justify-between font-mono">
          <div className="flex items-center gap-2">
            <CheckCircle2 className="w-4 h-4" />
            <span>Master Render Completed! Output URL: <strong>{renderResult.output_url}</strong></span>
          </div>
          <span className="bg-emerald-500/20 px-2 py-0.5 rounded text-emerald-300 font-medium">{renderResult.resolution}</span>
        </div>
      )}

      {/* Assembly Timeline UI */}
      <div className="bg-studio-card border border-studio-border rounded-md p-6 space-y-6">
        <h2 className="text-sm font-bold text-white uppercase tracking-wider font-mono">Episode 01 Multi-Track Assembly Timeline</h2>

        <div className="space-y-3 text-xs">
          {/* Video Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Video Track (9:16 Vertical Shading)</span>
            <div className="h-12 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto">
              <div className="h-full bg-studio-accent/20 border border-studio-accent/40 rounded px-3 flex items-center text-studio-text font-mono min-w-[160px]">
                [SHOT_001 - 00:00 - 00:05]
              </div>
              <div className="h-full bg-studio-accent/20 border border-studio-accent/40 rounded px-3 flex items-center text-studio-text font-mono min-w-[160px]">
                [SHOT_002 - 00:05 - 00:10]
              </div>
            </div>
          </div>

          {/* Dialogue Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Dialogue Track</span>
            <div className="h-10 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto">
              <div className="h-full bg-amber-500/10 border border-amber-500/30 rounded px-3 flex items-center text-amber-300 font-mono min-w-[220px]">
                [Sarah: "Did you find the file?"]
              </div>
            </div>
          </div>

          {/* Music & SFX Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Score & Ambient SFX Track</span>
            <div className="h-10 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto">
              <div className="h-full bg-emerald-500/10 border border-emerald-500/30 rounded px-3 flex items-center text-emerald-300 font-mono w-full">
                [Background Tension Score - Suspense Loop 120BPM]
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
