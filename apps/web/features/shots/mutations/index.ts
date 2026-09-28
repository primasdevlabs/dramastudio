import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Asset } from "@/lib/api/types";
import { shotKeys } from "../queries";

interface GenerateShotInput {
  project_id: string;
  shot_id: string;
  prompt: string;
  type: "video" | "image";
}

export function useGenerateShot() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: GenerateShotInput) =>
      api.post<Asset>("/v1/media/generate", input),
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
      shotId,
    }: {
      projectId: string;
      shotId: string;
    }) => api.post<Asset>(`/v1/media/assets/${shotId}/regenerate`, {}),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: shotKeys.all(variables.projectId),
      });
    },
  });
}
