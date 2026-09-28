import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { dialogueKeys } from "../queries";

// Dialogue generation goes through the generic AI generation endpoint
// with the dialogue_writing capability.
export function useGenerateDialogue() {
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
      api.post(`/v1/projects/${projectId}/ai/generate`, {
        capability: "dialogue_writing",
        input: { episode_id: episodeId },
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: dialogueKeys.all(
          variables.projectId,
          variables.seasonId,
          variables.episodeId
        ),
      });
    },
  });
}
