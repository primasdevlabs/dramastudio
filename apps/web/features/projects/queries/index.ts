import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Project } from "@/lib/api/types";

export const projectKeys = {
  all: ["projects"] as const,
  detail: (id: string) => ["projects", id] as const,
};

export function useProjects() {
  return useQuery({
    queryKey: projectKeys.all,
    queryFn: async () => {
      const res = await api.get<{ projects: Project[] }>("/v1/projects");
      return res.projects ?? [];
    },
  });
}

export function useProject(projectId: string) {
  return useQuery({
    queryKey: projectKeys.detail(projectId),
    queryFn: () => api.get<Project>(`/v1/projects/${projectId}`),
    enabled: Boolean(projectId),
  });
}
