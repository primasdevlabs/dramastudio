import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { audioKeys } from "../queries";
import type { Asset } from "@/lib/api/types";

export function useGenerateAudio() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      episodeId,
      capability = "voice",
      mediaType,
      prompt,
    }: {
      projectId: string;
      episodeId: string;
      capability?: "voice" | "music" | "sfx_generation";
      mediaType?: "voice" | "music" | "sfx";
      prompt?: string;
    }) =>
      api.post<{ job: unknown; asset: Asset }>(
        `/v1/projects/${projectId}/media/assets`,
        {
          capability,
          type: mediaType ?? "voice",
          episode_id: episodeId,
          prompt: prompt ?? "",
        }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: audioKeys.all(variables.projectId, variables.episodeId),
      });
    },
  });
}
