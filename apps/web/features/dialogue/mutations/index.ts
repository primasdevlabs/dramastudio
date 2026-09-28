import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { dialogueKeys } from "../queries";

export function useGenerateDialogue() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (episodeId: string) =>
      api.post(`/v1/episodes/${episodeId}/dialogue/generate`, {}),
    onSuccess: (_, episodeId) => {
      queryClient.invalidateQueries({ queryKey: dialogueKeys.all(episodeId) });
    },
  });
}
