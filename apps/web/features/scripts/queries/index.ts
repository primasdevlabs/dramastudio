import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Episode } from "@/lib/api/types";

// The script is the episode's `script` field — no separate resource.
export const scriptKeys = {
  detail: (projectId: string, seasonId: string, episodeId: string) =>
    ["script", projectId, seasonId, episodeId] as const,
};

export function useScript(projectId: string, seasonId: string, episodeId: string) {
  return useQuery({
    queryKey: scriptKeys.detail(projectId, seasonId, episodeId),
    queryFn: () =>
      api
        .get<Episode>(
          `/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}`
        )
        .catch(() => null),
    enabled: Boolean(projectId && seasonId && episodeId),
  });
}
