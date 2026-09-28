import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface ProductionAnalytics {
  total_generations: number;
  total_cost: number;
  episodes_completed: number;
  continuity_issues_resolved: number;
  active_agents: number;
}

export const analyticsKeys = {
  summary: (projectId: string) => ["analytics", projectId] as const,
};

export function useProductionAnalytics(projectId: string) {
  return useQuery({
    queryKey: analyticsKeys.summary(projectId),
    queryFn: () =>
      api
        .get<ProductionAnalytics>(
          `/v1/projects/${projectId}/analytics`
        )
        .catch(() => null),
    enabled: Boolean(projectId),
  });
}
