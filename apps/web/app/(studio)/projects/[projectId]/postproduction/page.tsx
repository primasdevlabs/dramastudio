"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import { Sliders, CheckCircle2, Cpu } from "lucide-react";
import { api } from "@/lib/api/client";
import { RenderTask, Episode, Season, EpisodeTimeline } from "@/lib/api/types";

export default function PostproductionPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const [renderResult, setRenderResult] = useState<RenderTask | null>(null);
  const [episodeId, setEpisodeId] = useState("");

  const { data: seasons } = useQuery({
    queryKey: ["seasons", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: Season[] }>(`/v1/projects/${projectId}/seasons`);
      return res.items || [];
    },
    enabled: Boolean(projectId),
  });

  const { data: episodes } = useQuery({
    queryKey: ["project-episodes", projectId, (seasons || []).map((s) => s.id).join(",")],
    queryFn: async () => {
      const all: Episode[] = [];
      for (const s of seasons || []) {
        const res = await api.get<{ items: Episode[] }>(
          `/v1/projects/${projectId}/seasons/${s.id}/episodes`
        );
        all.push(...(res.items || []));
      }
      return all;
    },
    enabled: Boolean(projectId && seasons && seasons.length > 0),
  });

  useEffect(() => {
    if (!episodeId && episodes && episodes.length > 0) {
      setEpisodeId(episodes[0].id);
    }
  }, [episodes, episodeId]);

  const { data: timeline } = useQuery({
    queryKey: ["timeline", projectId, episodeId],
    queryFn: () =>
      api
        .get<EpisodeTimeline>(
          `/v1/projects/${projectId}/postproduction/timelines?episode_id=${episodeId}`
        )
        .catch(() => null),
    enabled: Boolean(projectId && episodeId),
  });

  const renderMutation = useMutation({
    mutationFn: async () => {
      return api.post<RenderTask>(`/v1/projects/${projectId}/postproduction/renders`, {
        episode_id: episodeId,
        format: "mp4",
        resolution: "1080x1920",
      });
    },
    onSuccess: (data) => {
      setRenderResult(data);
    },
  });

  const selected = (episodes || []).find((e) => e.id === episodeId);

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

        <div className="flex items-center gap-3">
          <select
            value={episodeId}
            onChange={(e) => setEpisodeId(e.target.value)}
            className="bg-studio-panel border border-studio-border rounded px-3 py-2 text-xs text-white"
          >
            {(episodes || []).map((e) => (
              <option key={e.id} value={e.id}>
                E{String(e.number).padStart(2, "0")} — {e.title}
              </option>
            ))}
            {(!episodes || episodes.length === 0) && <option value="">No episodes</option>}
          </select>
          <button
            onClick={() => renderMutation.mutate()}
            disabled={renderMutation.isPending || !episodeId}
            className="inline-flex items-center gap-2 bg-studio-accent hover:bg-studio-accent-dark disabled:opacity-50 text-white font-medium text-xs px-4 py-2.5 rounded transition-colors"
          >
            <Cpu className="w-4 h-4" />
            {renderMutation.isPending ? "Rendering master..." : "Render Master (FFmpeg 9:16)"}
          </button>
        </div>
      </div>

      {renderResult && (
        <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 p-4 rounded text-xs flex items-center justify-between font-mono">
          <div className="flex items-center gap-2">
            <CheckCircle2 className="w-4 h-4" />
            <span>Render {renderResult.status}! Output URL: <strong>{renderResult.output_url || "pending"}</strong></span>
          </div>
          <span className="bg-emerald-500/20 px-2 py-0.5 rounded text-emerald-300 font-medium">{renderResult.resolution}</span>
        </div>
      )}

      {/* Assembly Timeline UI */}
      <div className="bg-studio-card border border-studio-border rounded-md p-6 space-y-6">
        <h2 className="text-sm font-bold text-white uppercase tracking-wider font-mono">
          {selected ? `E${String(selected.number).padStart(2, "0")} ` : ""}Multi-Track Assembly Timeline
          {timeline ? ` — v${timeline.version} (${timeline.status})` : ""}
        </h2>

        <div className="space-y-3 text-xs">
          {/* Video Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Video Track (9:16 Vertical)</span>
            <div className="h-12 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto">
              {timeline && timeline.video_tracks.length > 0 ? (
                timeline.video_tracks.map((t) => (
                  <div key={t.id} className="h-full bg-studio-accent/20 border border-studio-accent/40 rounded px-3 flex items-center text-studio-text font-mono min-w-[160px]">
                    [{t.shot_id} — {t.start_time}s +{t.duration}s]
                  </div>
                ))
              ) : (
                <div className="h-full flex items-center text-studio-muted font-mono px-2">No video track items yet</div>
              )}
            </div>
          </div>

          {/* Audio Track */}
          <div className="space-y-1">
            <span className="text-studio-muted font-semibold">Audio Track (dialogue / score / SFX)</span>
            <div className="h-10 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto">
              {timeline && timeline.audio_tracks.length > 0 ? (
                timeline.audio_tracks.map((t) => (
                  <div key={t.id} className="h-full bg-amber-500/10 border border-amber-500/30 rounded px-3 flex items-center text-amber-300 font-mono min-w-[220px]">
                    [{t.id} — {t.start_time}s +{t.duration}s]
                  </div>
                ))
              ) : (
                <div className="h-full flex items-center text-studio-muted font-mono px-2">No audio track items yet</div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
