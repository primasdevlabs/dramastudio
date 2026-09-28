import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { storyboardKeys, type Storyboard } from "../queries";

export function useGenerateStoryboard() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (episodeId: string) =>
      api.post<Storyboard>(`/v1/episodes/${episodeId}/storyboard/generate`, {}),
    onSuccess: (_, episodeId) => {
      queryClient.invalidateQueries({
        queryKey: storyboardKeys.detail(episodeId),
      });
    },
  });
}

export function useApproveStoryboard() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (episodeId: string) =>
      api.post(`/v1/episodes/${episodeId}/storyboard/approve`, {}),
    onSuccess: (_, episodeId) => {
      queryClient.invalidateQueries({
        queryKey: storyboardKeys.detail(episodeId),
      });
    },
  });
}
