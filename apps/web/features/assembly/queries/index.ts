import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { EpisodeTimeline } from "@/lib/api/types";

export const assemblyKeys = {
  detail: (projectId: string, episodeId: string) =>
    ["assembly", projectId, episodeId] as const,
};

export function useAssembly(projectId: string, episodeId: string) {
  return useQuery({
    queryKey: assemblyKeys.detail(projectId, episodeId),
    queryFn: () =>
      api
        .get<EpisodeTimeline>(
          `/v1/projects/${projectId}/postproduction/timelines?episode_id=${episodeId}`
        )
        .catch(() => null),
    enabled: Boolean(projectId && episodeId),
  });
}
