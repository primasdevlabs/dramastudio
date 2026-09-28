import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Episode } from "@/lib/api/types";

export const episodeKeys = {
  all: (seasonId: string) => ["episodes", seasonId] as const,
  detail: (episodeId: string) => ["episodes", "detail", episodeId] as const,
};

export function useEpisodes(seasonId: string) {
  return useQuery({
    queryKey: episodeKeys.all(seasonId),
    queryFn: async () => {
      const res = await api.get<{ episodes: Episode[] }>(
        `/v1/seasons/${seasonId}/episodes`
      );
      return res.episodes ?? [];
    },
    enabled: Boolean(seasonId),
  });
}

export function useEpisode(episodeId: string) {
  return useQuery({
    queryKey: episodeKeys.detail(episodeId),
    queryFn: () => api.get<Episode>(`/v1/episodes/${episodeId}`),
    enabled: Boolean(episodeId),
  });
}
