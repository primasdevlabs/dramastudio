"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Film, Sparkles, Play, Camera, CheckCircle2, Clock } from "lucide-react";
import { api } from "@/lib/api/client";
import { Asset } from "@/lib/api/types";

export default function StoryboardPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [prompt, setPrompt] = useState("Sarah looking suspiciously at phone in dark apartment");
  const [shotId, setShotId] = useState("shot_001");
  const [provider, setProvider] = useState("wan");

  const { data: assets, isLoading } = useQuery({
    queryKey: ["assets", projectId],
    queryFn: async () => {
      const res = await api.get<{ assets: Asset[] }>(`/v1/media/assets?project_id=${projectId}`);
      return res.assets || [];
    },
  });

  const generateMutation = useMutation({
    mutationFn: async () => {
      return api.post<Asset>("/v1/media/generate", {
        project_id: projectId,
        shot_id: shotId,
        prompt,
        provider,
        type: "video",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assets", projectId] });
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-indigo-500/10 border border-indigo-500/30 flex items-center justify-center text-indigo-400">
            <Film className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Storyboard & Shot Visualizer</h1>
            <p className="text-xs text-studio-muted">Plan camera specs, 9:16 vertical video shots, and trigger Wan Video Model provider generations</p>
          </div>
        </div>
      </div>

      {/* Generation Bar */}
      <div className="bg-studio-card border border-studio-border rounded-2xl p-5 space-y-4">
        <div className="flex items-center gap-2">
          <Sparkles className="w-4 h-4 text-cyan-400" />
          <h2 className="text-sm font-bold text-white">Trigger Wan Video Shot Generation</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="md:col-span-2">
            <label className="block text-xs font-semibold text-studio-muted mb-1">Structured Prompt / Shot Direction</label>
            <input
              type="text"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-cyan-400"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-studio-muted mb-1">Shot Identifier</label>
            <input
              type="text"
              value={shotId}
              onChange={(e) => setShotId(e.target.value)}
              className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-cyan-400"
            />
          </div>

          <div className="flex items-end">
            <button
              onClick={() => generateMutation.mutate()}
              disabled={generateMutation.isPending || !prompt}
              className="w-full inline-flex items-center justify-center gap-2 bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white font-bold text-xs py-2.5 rounded-lg shadow-lg shadow-cyan-500/20 transition-all"
            >
              <Camera className="w-4 h-4" />
              {generateMutation.isPending ? "Generating via Wan..." : "Generate Video Shot"}
            </button>
          </div>
        </div>
      </div>

      {/* Shot Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isLoading ? (
          [1, 2, 3].map((i) => <div key={i} className="h-64 bg-studio-card rounded-xl animate-pulse" />)
        ) : assets && assets.length > 0 ? (
          assets.map((asset) => (
            <div
              key={asset.id}
              className="bg-studio-card border border-studio-border rounded-xl overflow-hidden shadow-xl hover:border-cyan-500/40 transition-all flex flex-col justify-between"
            >
              <div className="relative aspect-[9/16] bg-black/80 flex items-center justify-center max-h-72 overflow-hidden">
                <div className="text-center p-4">
                  <div className="w-12 h-12 rounded-full bg-cyan-500/20 text-cyan-400 flex items-center justify-center mx-auto mb-2 border border-cyan-500/30">
                    <Play className="w-6 h-6 ml-0.5" />
                  </div>
                  <span className="text-xs font-mono text-cyan-400">9:16 Vertical Render Preview</span>
                </div>
              </div>

              <div className="p-4 space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-bold text-white uppercase tracking-wider">{asset.shot_id || "Shot"}</span>
                  <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 text-xs font-semibold border border-emerald-500/20">
                    {asset.provider.toUpperCase()} (Wan)
                  </span>
                </div>

                <p className="text-xs text-studio-muted line-clamp-2 bg-studio-panel p-2.5 rounded-lg border border-studio-border/60">
                  {asset.prompt}
                </p>

                <div className="flex items-center justify-between text-xs pt-2 border-t border-studio-border/60">
                  <span className="text-studio-muted">Version {asset.version}</span>
                  <span className="text-emerald-400 font-semibold">${asset.cost || 0.05}</span>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full bg-studio-card border border-studio-border rounded-2xl p-12 text-center text-studio-muted text-xs">
            No video shots generated yet. Use the prompt generator above to trigger Wan Video model generation.
          </div>
        )}
      </div>
    </div>
  );
}
