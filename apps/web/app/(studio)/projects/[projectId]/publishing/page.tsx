"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Share2, Send, CheckCircle2 } from "lucide-react";
import { api } from "@/lib/api/client";
import { Publication } from "@/lib/api/types";

export default function PublishingPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [channel, setChannel] = useState("tiktok_channel_1");
  const [title, setTitle] = useState("Episode 01: Midnight Call");
  const [caption, setCaption] = useState("She received an anonymous file. Watch the truth unfold! #drama #thriller");
  const [published, setPublished] = useState(false);

  const { data: publications } = useQuery({
    queryKey: ["publications", projectId],
    queryFn: async () => {
      const res = await api.get<{ publications: Publication[] }>(`/v1/publishing/publications?project_id=${projectId}`);
      return res.publications || [];
    },
  });

  const publishMutation = useMutation({
    mutationFn: async () => {
      return api.post<Publication>("/v1/publishing/publish", {
        episode_id: "ep_001",
        channel_id: channel,
        title,
        caption,
        tags: ["#drama", "#thriller", "#serialized"],
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["publications", projectId] });
      setPublished(true);
      setTimeout(() => setPublished(false), 3000);
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-md">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded bg-studio-accent/10 border border-studio-accent/20 flex items-center justify-center text-studio-accent">
            <Share2 className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Publishing Workspace</h1>
            <p className="text-xs text-studio-muted">Multi-channel distribution to TikTok, Instagram Reels, and YouTube Shorts</p>
          </div>
        </div>
      </div>

      {published && (
        <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 p-4 rounded text-xs font-mono flex items-center gap-2">
          <CheckCircle2 className="w-4 h-4" />
          Episode successfully published to {channel.toUpperCase()}!
        </div>
      )}

      {/* Publishing Form */}
      <div className="bg-studio-card border border-studio-border rounded-md p-6 space-y-4 max-w-2xl">
        <h2 className="text-base font-bold text-white">Publish Episode 01</h2>

        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Target Channel / Platform</label>
          <select
            value={channel}
            onChange={(e) => setChannel(e.target.value)}
            className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-xs text-white focus:outline-none focus:border-studio-accent font-mono"
          >
            <option value="tiktok_channel_1">TikTok Official Channel</option>
            <option value="instagram_reels_1">Instagram Reels Channel</option>
            <option value="youtube_shorts_1">YouTube Shorts Channel</option>
          </select>
        </div>

        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Publication Title</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-xs text-white focus:outline-none focus:border-studio-accent"
          />
        </div>

        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Caption & Hashtags</label>
          <textarea
            rows={3}
            value={caption}
            onChange={(e) => setCaption(e.target.value)}
            className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-xs text-white focus:outline-none focus:border-studio-accent resize-none"
          />
        </div>

        <button
          onClick={() => publishMutation.mutate()}
          disabled={publishMutation.isPending}
          className="inline-flex items-center justify-center gap-2 bg-studio-accent hover:bg-studio-accent-dark text-white font-medium text-xs px-6 py-2.5 rounded transition-colors"
        >
          <Send className="w-4 h-4" />
          {publishMutation.isPending ? "Publishing..." : "Publish Episode Now"}
        </button>
      </div>

      {/* Publications History Table */}
      <div className="bg-studio-card border border-studio-border rounded-md p-5 space-y-3">
        <h3 className="text-sm font-bold text-white">Publication History</h3>
        {publications && publications.length > 0 ? (
          <div className="space-y-2">
            {publications.map((p) => (
              <div key={p.id} className="bg-studio-panel border border-studio-border p-3.5 rounded flex items-center justify-between text-xs">
                <div>
                  <span className="font-bold text-white">{p.metadata?.title}</span>
                  <p className="text-studio-muted mt-0.5">{p.metadata?.caption}</p>
                </div>
                <span className="px-2.5 py-0.5 rounded bg-studio-panel text-studio-accent font-mono border border-studio-border">
                  {p.channel_id}
                </span>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-xs text-studio-muted italic">No publications recorded yet.</p>
        )}
      </div>
    </div>
  );
}
