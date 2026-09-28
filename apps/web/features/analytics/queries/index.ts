import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { MetricSummary } from "@/lib/api/types";

export interface Metric {
  id: number;
  project_id: string;
  episode_id?: string;
  publication_id?: string;
  metric: string;
  value: number;
  dimensions?: Record<string, unknown>;
  recorded_at: string;
}

export const analyticsKeys = {
  metrics: (projectId: string) => ["analytics-metrics", projectId] as const,
  summary: (projectId: string, metric: string) =>
    ["analytics-summary", projectId, metric] as const,
};

export function useMetrics(projectId: string) {
  return useQuery({
    queryKey: analyticsKeys.metrics(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: Metric[] }>(
        `/v1/projects/${projectId}/analytics/metrics`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useMetricSummary(projectId: string, metric: string) {
  return useQuery({
    queryKey: analyticsKeys.summary(projectId, metric),
    queryFn: () =>
      api
        .get<MetricSummary>(
          `/v1/projects/${projectId}/analytics/summary?metric=${metric}`
        )
        .catch(() => null),
    enabled: Boolean(projectId && metric),
  });
}
