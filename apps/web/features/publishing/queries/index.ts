import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Publication, Channel } from "@/lib/api/types";

export const publishingKeys = {
  all: (projectId: string) => ["publications", projectId] as const,
  channels: (projectId: string) => ["channels", projectId] as const,
};

export function usePublications(projectId: string) {
  return useQuery({
    queryKey: publishingKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: Publication[] }>(
        `/v1/projects/${projectId}/publishing/publications`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useChannels(projectId: string) {
  return useQuery({
    queryKey: publishingKeys.channels(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: Channel[] }>(
        `/v1/projects/${projectId}/publishing/channels`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
