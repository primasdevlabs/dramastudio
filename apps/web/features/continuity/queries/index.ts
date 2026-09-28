import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { ContinuityIssue } from "@/lib/api/types";

export const continuityKeys = {
  all: (projectId: string) => ["continuity-issues", projectId] as const,
};

export function useContinuityIssues(projectId: string) {
  return useQuery({
    queryKey: continuityKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ issues: ContinuityIssue[] }>(
        `/v1/continuity/issues?project_id=${projectId}`
      );
      return res.issues ?? [];
    },
    enabled: Boolean(projectId),
  });
}
