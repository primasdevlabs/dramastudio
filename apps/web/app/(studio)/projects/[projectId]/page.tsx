"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import {
  Clapperboard,
  BookOpen,
  Users,
  Globe,
  Film,
  Activity,
  ShieldCheck,
  ChevronRight,
  Sparkles,
  Layers,
  Sliders,
  Share2,
} from "lucide-react";
import { api } from "@/lib/api/client";
import { Project, SeriesBible } from "@/lib/api/types";

export default function ProjectDashboardPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;

  const { data: project } = useQuery({
    queryKey: ["project", projectId],
    queryFn: () => api.get<Project>(`/v1/projects/${projectId}`),
  });

  const { data: bible } = useQuery({
    queryKey: ["bible", projectId],
    queryFn: () => api.get<SeriesBible>(`/v1/projects/${projectId}/bible`).catch(() => null),
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Studio Header Banner */}
      <div className="bg-studio-card border border-studio-border p-6 rounded-md flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-2 font-mono">
            <span className="px-2 py-0.5 rounded bg-studio-accent/10 text-studio-accent text-xs font-medium border border-studio-accent/20">
              {project?.genre || "Drama"}
            </span>
            <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 text-xs font-medium border border-emerald-500/20">
              {project?.status || "SERIES_APPROVED"}
            </span>
          </div>
          <h1 className="text-2xl font-bold text-white tracking-tight">{project?.name || "Loading Project..."}</h1>
          <p className="text-xs text-studio-muted mt-1">{project?.description || "Serialized AI Drama Production Workspace"}</p>
        </div>

        <div className="flex items-center gap-3">
          <Link
            href={`/projects/${projectId}/production`}
            className="inline-flex items-center gap-2 bg-studio-accent hover:bg-studio-accent-dark text-white font-medium text-xs px-4 py-2 rounded transition-colors"
          >
            <Activity className="w-4 h-4" /> Production Tower
          </Link>
        </div>
      </div>

      {/* Production Quick Navigation Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          { title: "Series Bible", desc: "Premise, tone, world & narrative rules", path: "/bible", icon: BookOpen, color: "text-amber-400" },
          { title: "Story Workspace", desc: "Seasons, episodes, arcs & story graph", path: "/story", icon: Layers, color: "text-studio-accent" },
          { title: "Character Studio", desc: "Identity, wardrobe, voice profiles", path: "/characters", icon: Users, color: "text-emerald-400" },
          { title: "World Studio", desc: "Locations, environments & props", path: "/world", icon: Globe, color: "text-emerald-400" },
          { title: "Storyboards & Shots", path: "/storyboard", desc: "Camera angles, shots & 9:16 visuals", icon: Film, color: "text-amber-400" },
          { title: "Continuity Center", desc: "Wardrobe & story timeline checkers", path: "/continuity", icon: ShieldCheck, color: "text-rose-400" },
          { title: "Postproduction", desc: "Assembly timeline, FFmpeg renderer", path: "/postproduction", icon: Sliders, color: "text-studio-accent" },
          { title: "Publishing", desc: "TikTok, Instagram & YouTube channels", path: "/publishing", icon: Share2, color: "text-studio-muted" },
        ].map((item) => {
          const Icon = item.icon;
          return (
            <Link
              key={item.title}
              href={`/projects/${projectId}${item.path}`}
              className="bg-studio-card border border-studio-border p-4 rounded-md hover:border-studio-accent/60 transition-all group"
            >
              <div className="flex items-center justify-between mb-2">
                <div className={`p-2 rounded bg-studio-panel border border-studio-border ${item.color}`}>
                  <Icon className="w-5 h-5" />
                </div>
                <ChevronRight className="w-4 h-4 text-studio-muted group-hover:text-studio-accent transition-colors" />
              </div>
              <h3 className="text-sm font-bold text-white group-hover:text-studio-accent transition-colors">{item.title}</h3>
              <p className="text-xs text-studio-muted mt-1">{item.desc}</p>
            </Link>
          );
        })}
      </div>

      {/* Series Bible Snapshot Card */}
      <div className="bg-studio-card border border-studio-border rounded-md p-6 space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <BookOpen className="w-5 h-5 text-amber-400" />
            <h2 className="text-lg font-bold text-white">Active Series Bible</h2>
            {bible && (
              <span className="text-xs font-mono px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20">
                v{bible.version}
              </span>
            )}
          </div>
          <Link
            href={`/projects/${projectId}/bible`}
            className="text-xs font-semibold text-studio-accent hover:underline"
          >
            Edit Bible
          </Link>
        </div>

        {bible ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
            <div className="bg-studio-panel border border-studio-border p-3.5 rounded space-y-1">
              <span className="text-studio-muted font-semibold">Premise</span>
              <p className="text-studio-text">{bible.premise}</p>
            </div>
            <div className="bg-studio-panel border border-studio-border p-3.5 rounded space-y-1">
              <span className="text-studio-muted font-semibold">Themes</span>
              <div className="flex flex-wrap gap-1 mt-1 font-mono">
                {bible.themes?.map((t) => (
                  <span key={t} className="px-2 py-0.5 rounded bg-studio-panel text-studio-text border border-studio-border">
                    {t}
                  </span>
                ))}
              </div>
            </div>
          </div>
        ) : (
          <div className="text-xs text-studio-muted italic">
            No Series Bible created yet. <Link href={`/projects/${projectId}/bible`} className="text-studio-accent underline">Initialize Series Bible</Link>
          </div>
        )}
      </div>
    </div>
  );
}
