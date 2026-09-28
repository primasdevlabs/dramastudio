import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Asset } from "@/lib/api/types";
import { shotKeys } from "../queries";

interface GenerateShotInput {
  project_id: string;
  episode_id?: string;
  shot_id?: string;
  prompt: string;
  type: "video" | "image";
}

export function useGenerateShot() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ project_id, type, ...rest }: GenerateShotInput) =>
      api.post<{ job: unknown; asset: Asset }>(
        `/v1/projects/${project_id}/media/assets`,
        {
          capability: type === "video" ? "video_generation" : "image_generation",
          type,
          ...rest,
        }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: shotKeys.all(variables.project_id),
      });
    },
  });
}

export function useRegenerateShot() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      assetId,
    }: {
      projectId: string;
      assetId: string;
    }) =>
      api.post<Asset>(
        `/v1/projects/${projectId}/media/assets/${assetId}/regenerate`,
        {}
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: shotKeys.all(variables.projectId),
      });
    },
  });
}
