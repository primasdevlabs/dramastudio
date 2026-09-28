import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Episode } from "@/lib/api/types";

export const episodeKeys = {
  all: (projectId: string, seasonId: string) =>
    ["episodes", projectId, seasonId] as const,
  detail: (projectId: string, seasonId: string, episodeId: string) =>
    ["episodes", "detail", projectId, seasonId, episodeId] as const,
};

export function useEpisodes(projectId: string, seasonId: string) {
  return useQuery({
    queryKey: episodeKeys.all(projectId, seasonId),
    queryFn: async () => {
      const res = await api.get<{ items: Episode[] }>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId && seasonId),
  });
}

export function useEpisode(projectId: string, seasonId: string, episodeId: string) {
  return useQuery({
    queryKey: episodeKeys.detail(projectId, seasonId, episodeId),
    queryFn: () =>
      api.get<Episode>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}`
      ),
    enabled: Boolean(projectId && seasonId && episodeId),
  });
}
