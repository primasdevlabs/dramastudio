import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { LeadDirectorDecision, ApprovalRequest } from "@/lib/api/types";

export const productionKeys = {
  decisions: (projectId: string) => ["decisions", projectId] as const,
  approvals: (projectId: string) => ["approvals", projectId] as const,
};

export function useDecisions(projectId: string) {
  return useQuery({
    queryKey: productionKeys.decisions(projectId),
    queryFn: async () => {
      const res = await api.get<{ decisions: LeadDirectorDecision[] }>(
        `/v1/agents/decisions?project_id=${projectId}`
      );
      return res.decisions ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useApprovals(projectId: string) {
  return useQuery({
    queryKey: productionKeys.approvals(projectId),
    queryFn: async () => {
      const res = await api.get<{ approvals: ApprovalRequest[] }>(
        `/v1/production/approval?project_id=${projectId}`
      );
      return res.approvals ?? [];
    },
    enabled: Boolean(projectId),
  });
}
