import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Episode } from "@/lib/api/types";
import { episodeKeys } from "../queries";

interface CreateEpisodeInput {
  season_id: string;
  number: number;
  title: string;
  summary: string;
}

export function useCreateEpisode() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ season_id, ...data }: CreateEpisodeInput) =>
      api.post<Episode>(`/v1/seasons/${season_id}/episodes`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: episodeKeys.all(variables.season_id),
      });
    },
  });
}
