import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Publication } from "@/lib/api/types";
import { publishingKeys } from "../queries";

interface PublishInput {
  project_id: string;
  episode_id: string;
  channel_id: string;
  title: string;
  caption: string;
  tags: string[];
}

export function usePublishEpisode() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ project_id, ...data }: PublishInput) =>
      api.post<Publication>("/v1/publishing/publish", data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: publishingKeys.all(variables.project_id),
      });
    },
  });
}
