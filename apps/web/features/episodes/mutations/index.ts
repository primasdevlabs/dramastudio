import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Episode } from "@/lib/api/types";
import { episodeKeys } from "../queries";

interface CreateEpisodeInput {
  project_id: string;
  season_id: string;
  number: number;
  title: string;
  summary: string;
}

export function useCreateEpisode() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ project_id, season_id, ...data }: CreateEpisodeInput) =>
      api.post<Episode>(
        `/v1/projects/${project_id}/seasons/${season_id}/episodes`,
        data
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: episodeKeys.all(variables.project_id, variables.season_id),
      });
    },
  });
}
