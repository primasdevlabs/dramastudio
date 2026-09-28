import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Season } from "@/lib/api/types";

export const seasonKeys = {
  all: (projectId: string) => ["seasons", projectId] as const,
  detail: (projectId: string, seasonId: string) =>
    ["seasons", "detail", projectId, seasonId] as const,
};

export function useSeasons(projectId: string) {
  return useQuery({
    queryKey: seasonKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: Season[] }>(
        `/v1/projects/${projectId}/seasons`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useSeason(projectId: string, seasonId: string) {
  return useQuery({
    queryKey: seasonKeys.detail(projectId, seasonId),
    queryFn: () =>
      api.get<Season>(`/v1/projects/${projectId}/seasons/${seasonId}`),
    enabled: Boolean(projectId && seasonId),
  });
}
