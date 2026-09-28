import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Publication } from "@/lib/api/types";
import { publishingKeys } from "../queries";

interface PublishInput {
  project_id: string;
  episode_id: string;
  channel_id: string;
  video_url?: string;
  title: string;
  caption: string;
  tags: string[];
  scheduled_at?: string;
}

export function usePublishEpisode() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ project_id, title, caption, tags, ...data }: PublishInput) =>
      api.post<Publication>(
        `/v1/projects/${project_id}/publishing/publications`,
        {
          ...data,
          metadata: { title, caption, tags },
        }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: publishingKeys.all(variables.project_id),
      });
    },
  });
}

export function usePublishNow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      publicationId,
    }: {
      projectId: string;
      publicationId: string;
    }) =>
      api.post<Publication>(
        `/v1/projects/${projectId}/publishing/publications/${publicationId}/publish`,
        {}
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: publishingKeys.all(variables.projectId),
      });
    },
  });
}
