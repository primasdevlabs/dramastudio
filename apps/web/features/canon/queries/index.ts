import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { StoryFact } from "@/lib/api/types";

export const canonKeys = {
  all: (projectId: string) => ["canon-facts", projectId] as const,
};

export function useCanonFacts(projectId: string) {
  return useQuery({
    queryKey: canonKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ facts: StoryFact[] }>(
        `/v1/canon/facts?project_id=${projectId}`
      );
      return res.facts ?? [];
    },
    enabled: Boolean(projectId),
  });
}
