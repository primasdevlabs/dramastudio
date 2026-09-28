import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Season } from "@/lib/api/types";

export const seasonKeys = {
  all: (projectId: string) => ["seasons", projectId] as const,
  detail: (seasonId: string) => ["seasons", "detail", seasonId] as const,
};

export function useSeasons(projectId: string) {
  return useQuery({
    queryKey: seasonKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ seasons: Season[] }>(
        `/v1/projects/${projectId}/seasons`
      );
      return res.seasons ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useSeason(seasonId: string) {
  return useQuery({
    queryKey: seasonKeys.detail(seasonId),
    queryFn: () => api.get<Season>(`/v1/seasons/${seasonId}`),
    enabled: Boolean(seasonId),
  });
}
