import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Scene } from "@/lib/api/types";

// Dialogue lives on scenes (action/dialogue fields) — there is no
// standalone dialogue resource in the backend.
export const dialogueKeys = {
  all: (projectId: string, seasonId: string, episodeId: string) =>
    ["dialogue", projectId, seasonId, episodeId] as const,
};

export function useDialogue(projectId: string, seasonId: string, episodeId: string) {
  return useQuery({
    queryKey: dialogueKeys.all(projectId, seasonId, episodeId),
    queryFn: async () => {
      const res = await api.get<{ items: Scene[] }>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}/scenes`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId && seasonId && episodeId),
  });
}
