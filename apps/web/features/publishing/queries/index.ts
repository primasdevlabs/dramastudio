import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Publication } from "@/lib/api/types";

export const publishingKeys = {
  all: (projectId: string) => ["publications", projectId] as const,
};

export function usePublications(projectId: string) {
  return useQuery({
    queryKey: publishingKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ publications: Publication[] }>(
        `/v1/publishing/publications?project_id=${projectId}`
      );
      return res.publications ?? [];
    },
    enabled: Boolean(projectId),
  });
}
