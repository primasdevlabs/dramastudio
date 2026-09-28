import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface AgentTask {
  id: string;
  agent_id: string;
  objective: string;
  status: string;
  result?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export const agentKeys = {
  tasks: (projectId: string) => ["agent-tasks", projectId] as const,
};

export function useAgentTasks(projectId: string) {
  return useQuery({
    queryKey: agentKeys.tasks(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: AgentTask[] }>(
        `/v1/projects/${projectId}/agents/tasks`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
