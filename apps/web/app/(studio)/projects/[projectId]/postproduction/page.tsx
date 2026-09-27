"use client";

import { useState } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import { Sliders, Play, Film, CheckCircle2, Cpu } from "lucide-react";
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
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-violet-500/10 border border-violet-500/30 flex items-center justify-center text-violet-400">
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
          className="inline-flex items-center gap-2 bg-gradient-to-r from-violet-500 to-indigo-600 hover:from-violet-400 hover:to-indigo-500 text-white font-bold text-xs px-4 py-2.5 rounded-lg shadow-lg shadow-violet-500/20 transition-all"
        >
          <Cpu className="w-4 h-4" />
          {renderMutation.isPending ? "Rendering via FFmpeg..." : "Render Episode (FFmpeg 9:16)"}
        </button>
      </div>

      {renderResult && (
        <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 p-4 rounded-2xl text-xs flex items-center justify-between">
          <div className="flex items-center gap-2">
            <CheckCircle2 className="w-4 h-4" />
            <span>FFmpeg Render Task Completed! Output URL: <strong>{renderResult.output_url}</strong></span>
          </div>
          <span className="font-mono bg-emerald-500/20 px-2 py-0.5 rounded">{renderResult.resolution}</span>
        </div>
      )}

      {/* Assembly Timeline UI */}
      <div className="bg-studio-card border border-studio-border rounded-2xl p-6 space-y-6">
        <h2 className="text-sm font-bold text-white uppercase tracking-wider">Episode 01 Multi-Track Assembly Timeline</h2>

        <div className="space-y-3 text-xs">
          {/* Video Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Video Track (9:16 Vertical Shading)</span>
            <div className="h-12 bg-studio-panel border border-studio-border rounded-xl p-2 flex gap-2 overflow-x-auto">
              <div className="h-full bg-cyan-500/20 border border-cyan-500/40 rounded px-3 flex items-center font-mono text-cyan-300 min-w-[140px]">
                [Shot 001 - 0:00 - 0:05]
              </div>
              <div className="h-full bg-indigo-500/20 border border-indigo-500/40 rounded px-3 flex items-center font-mono text-indigo-300 min-w-[140px]">
                [Shot 002 - 0:05 - 0:10]
              </div>
            </div>
          </div>

          {/* Dialogue Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Dialogue Track</span>
            <div className="h-10 bg-studio-panel border border-studio-border rounded-xl p-2 flex gap-2 overflow-x-auto">
              <div className="h-full bg-pink-500/20 border border-pink-500/40 rounded px-3 flex items-center font-mono text-pink-300 min-w-[200px]">
                [Sarah: "Did you find the file?"]
              </div>
            </div>
          </div>

          {/* Music & SFX Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Score & Ambient SFX Track</span>
            <div className="h-10 bg-studio-panel border border-studio-border rounded-xl p-2 flex gap-2 overflow-x-auto">
              <div className="h-full bg-emerald-500/20 border border-emerald-500/40 rounded px-3 flex items-center font-mono text-emerald-300 w-full">
                [Background Tension Score - Suspense Loop 120BPM]
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
