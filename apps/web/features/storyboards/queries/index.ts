import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Shot } from "@/lib/api/types";

// The storyboard is the episode's ordered shot list in production.
export const storyboardKeys = {
  shots: (projectId: string, episodeId: string) =>
    ["storyboard-shots", projectId, episodeId] as const,
};

export function useStoryboard(projectId: string, episodeId: string) {
  return useQuery({
    queryKey: storyboardKeys.shots(projectId, episodeId),
    queryFn: async () => {
      const res = await api.get<{ items: Shot[] }>(
        `/v1/projects/${projectId}/production/shots?episode_id=${episodeId}`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId && episodeId),
  });
}
