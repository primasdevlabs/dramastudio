import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Episode } from "@/lib/api/types";
import { scriptKeys } from "../queries";

export function useGenerateScript() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      episodeId,
    }: {
      projectId: string;
      seasonId: string;
      episodeId: string;
    }) =>
      api.post(`/v1/projects/${projectId}/ai/generate`, {
        capability: "script_writing",
        input: { episode_id: episodeId },
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: scriptKeys.detail(
          variables.projectId,
          variables.seasonId,
          variables.episodeId
        ),
      });
    },
  });
}

// Approving a script transitions the episode to SCRIPTED.
export function useApproveScript() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      seasonId,
      episodeId,
    }: {
      projectId: string;
      seasonId: string;
      episodeId: string;
    }) =>
      api.patch<Episode>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}`,
        { status: "SCRIPTED" }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: scriptKeys.detail(
          variables.projectId,
          variables.seasonId,
          variables.episodeId
        ),
      });
    },
  });
}
