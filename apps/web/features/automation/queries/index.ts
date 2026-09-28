import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface WorkflowStatus {
  id: string;
  name: string;
  status: "RUNNING" | "PAUSED" | "COMPLETED" | "FAILED";
  current_step: string;
  started_at: string;
}

export const automationKeys = {
  workflows: (projectId: string) => ["workflows", projectId] as const,
};

export function useWorkflows(projectId: string) {
  return useQuery({
    queryKey: automationKeys.workflows(projectId),
    queryFn: async () => {
      const res = await api.get<{ workflows: WorkflowStatus[] }>(
        `/v1/automation/workflows?project_id=${projectId}`
      );
      return res.workflows ?? [];
    },
    enabled: Boolean(projectId),
  });
}
