import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Asset } from "@/lib/api/types";

export const audioKeys = {
  all: (projectId: string, episodeId: string) =>
    ["audio", projectId, episodeId] as const,
};

export function useAudioTracks(projectId: string, episodeId: string) {
  return useQuery({
    queryKey: audioKeys.all(projectId, episodeId),
    queryFn: async () => {
      // The asset list is project-scoped; episode + audio-type filtering
      // happens client-side (backend has no episode_id query param).
      const res = await api.get<{ items: Asset[] }>(
        `/v1/projects/${projectId}/media/assets`
      );
      return (res.items ?? []).filter(
        (a) =>
          a.episode_id === episodeId &&
          (a.type === "voice" || a.type === "music" || a.type === "sfx")
      );
    },
    enabled: Boolean(projectId && episodeId),
  });
}
