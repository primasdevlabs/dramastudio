import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { LeadDirectorDecision, ApprovalRequest, ProductionRun } from "@/lib/api/types";

export const productionKeys = {
  decisions: (projectId: string) => ["decisions", projectId] as const,
  approvals: (projectId: string) => ["approvals", projectId] as const,
  runs: (projectId: string) => ["production-runs", projectId] as const,
};

export function useDecisions(projectId: string) {
  return useQuery({
    queryKey: productionKeys.decisions(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: LeadDirectorDecision[] }>(
        `/v1/projects/${projectId}/agents/decisions`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useApprovals(projectId: string) {
  return useQuery({
    queryKey: productionKeys.approvals(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: ApprovalRequest[] }>(
        `/v1/projects/${projectId}/production/approvals`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useProductionRuns(projectId: string) {
  return useQuery({
    queryKey: productionKeys.runs(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: ProductionRun[] }>(
        `/v1/projects/${projectId}/production/runs`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
