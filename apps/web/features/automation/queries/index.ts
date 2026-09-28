import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { ProductionRun } from "@/lib/api/types";

export const automationKeys = {
  runs: (projectId: string) => ["production-runs", projectId] as const,
};

// Automation workflows are realized as production runs: the backend
// state machine drives them, and pause/resume/stop map to run actions.
export function useWorkflows(projectId: string) {
  return useQuery({
    queryKey: automationKeys.runs(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: ProductionRun[] }>(
        `/v1/projects/${projectId}/production/runs`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
