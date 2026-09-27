"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { FolderPlus, Sparkles, Activity } from "lucide-react";
import { api } from "@/lib/api/client";
import { Project } from "@/lib/api/types";
import { ProjectCard } from "@/features/projects/components/project-card";
import { NewProjectDialog } from "@/features/projects/components/new-project-dialog";

export default function ProjectsPage() {
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["projects"],
    queryFn: async () => {
      const res = await api.get<{ projects: Project[] }>("/v1/projects");
      return res.projects || [];
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header Banner */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <span className="px-2.5 py-0.5 rounded-full bg-cyan-500/10 text-cyan-400 text-xs font-semibold border border-cyan-500/20">
              Control Tower
            </span>
          </div>
          <h1 className="text-2xl font-black text-white tracking-tight">Active Drama Productions</h1>
          <p className="text-xs text-studio-muted">
            Manage serialized AI drama projects, series bibles, and production pipelines
          </p>
        </div>
        <NewProjectDialog onCreated={() => queryClient.invalidateQueries({ queryKey: ["projects"] })} />
      </div>

      {/* Grid List */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-44 bg-studio-card border border-studio-border rounded-xl animate-pulse" />
          ))}
        </div>
      ) : error ? (
        <div className="bg-rose-500/10 border border-rose-500/20 text-rose-400 p-4 rounded-xl text-sm">
          Failed to load projects: {(error as Error).message}
        </div>
      ) : data && data.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {data.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      ) : (
        <div className="bg-studio-card border border-studio-border rounded-2xl p-12 text-center max-w-md mx-auto space-y-4">
          <div className="w-12 h-12 rounded-2xl bg-cyan-500/10 border border-cyan-500/20 text-cyan-400 flex items-center justify-center mx-auto">
            <FolderPlus className="w-6 h-6" />
          </div>
          <div>
            <h3 className="text-lg font-bold text-white">No Productions Yet</h3>
            <p className="text-xs text-studio-muted mt-1">
              Start by initializing your first AI drama series project.
            </p>
          </div>
          <NewProjectDialog onCreated={() => queryClient.invalidateQueries({ queryKey: ["projects"] })} />
        </div>
      )}
    </div>
  );
}
