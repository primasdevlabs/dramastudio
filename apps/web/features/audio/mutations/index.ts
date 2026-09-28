import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { audioKeys } from "../queries";

export function useGenerateAudio() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (episodeId: string) =>
      api.post(`/v1/episodes/${episodeId}/audio/generate`, {}),
    onSuccess: (_, episodeId) => {
      queryClient.invalidateQueries({ queryKey: audioKeys.all(episodeId) });
    },
  });
}
