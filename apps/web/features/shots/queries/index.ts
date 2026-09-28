import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Shot } from "@/lib/api/types";

export const shotKeys = {
  all: (projectId: string, episodeId?: string) =>
    ["shots", projectId, episodeId ?? "all"] as const,
  detail: (projectId: string, shotId: string) =>
    ["shots", "detail", projectId, shotId] as const,
};

export function useShots(projectId: string, episodeId?: string) {
  return useQuery({
    queryKey: shotKeys.all(projectId, episodeId),
    queryFn: async () => {
      const qs = episodeId ? `?episode_id=${episodeId}` : "";
      const res = await api.get<{ items: Shot[] }>(
        `/v1/projects/${projectId}/production/shots${qs}`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
