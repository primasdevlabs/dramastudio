import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Shot } from "@/lib/api/types";
import { storyboardKeys } from "../queries";

// Storyboard generation is the storyboard AI capability; the resulting
// shots materialize in production.
export function useGenerateStoryboard() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      episodeId,
    }: {
      projectId: string;
      episodeId: string;
    }) =>
      api.post(`/v1/projects/${projectId}/ai/generate`, {
        capability: "storyboard",
        input: { episode_id: episodeId },
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: storyboardKeys.shots(variables.projectId, variables.episodeId),
      });
    },
  });
}

export function useApproveShot() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      shotId,
      episodeId,
    }: {
      projectId: string;
      shotId: string;
      episodeId: string;
    }) =>
      api.post<Shot>(
        `/v1/projects/${projectId}/production/shots/${shotId}/approve`,
        {}
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: storyboardKeys.shots(variables.projectId, variables.episodeId),
      });
    },
  });
}
