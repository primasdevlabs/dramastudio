import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { SeriesBible } from "@/lib/api/types";

export const bibleKeys = {
  detail: (projectId: string) => ["bible", projectId] as const,
};

export function useBible(projectId: string) {
  return useQuery({
    queryKey: bibleKeys.detail(projectId),
    queryFn: () =>
      api.get<SeriesBible>(`/v1/projects/${projectId}/bible`).catch(() => null),
    enabled: Boolean(projectId),
  });
}
