import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { scriptKeys, type Script } from "../queries";

export function useGenerateScript() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (episodeId: string) =>
      api.post<Script>(`/v1/episodes/${episodeId}/script/generate`, {}),
    onSuccess: (_, episodeId) => {
      queryClient.invalidateQueries({ queryKey: scriptKeys.detail(episodeId) });
    },
  });
}

export function useApproveScript() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (episodeId: string) =>
      api.post<Script>(`/v1/episodes/${episodeId}/script/approve`, {}),
    onSuccess: (_, episodeId) => {
      queryClient.invalidateQueries({ queryKey: scriptKeys.detail(episodeId) });
    },
  });
}
