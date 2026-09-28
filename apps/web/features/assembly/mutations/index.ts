import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { assemblyKeys } from "../queries";
import type { RenderTask } from "@/lib/api/types";

export function useRenderEpisode() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      episodeId,
      format,
      resolution,
    }: {
      projectId: string;
      episodeId: string;
      format: string;
      resolution: string;
    }) =>
      api.post<RenderTask>(
        `/v1/projects/${projectId}/postproduction/renders`,
        { episode_id: episodeId, format, resolution }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: assemblyKeys.detail(variables.projectId, variables.episodeId),
      });
    },
  });
}
