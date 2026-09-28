import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface AgentActivity {
  agent_name: string;
  role: string;
  current_task: string;
  status: "ACTIVE" | "IDLE" | "WAITING" | "ERROR";
  last_action_at: string;
}

export const agentKeys = {
  activity: (projectId: string) => ["agent-activity", projectId] as const,
};

export function useAgentActivity(projectId: string) {
  return useQuery({
    queryKey: agentKeys.activity(projectId),
    queryFn: async () => {
      const res = await api.get<{ agents: AgentActivity[] }>(
        `/v1/agents/activity?project_id=${projectId}`
      );
      return res.agents ?? [];
    },
    enabled: Boolean(projectId),
  });
}
