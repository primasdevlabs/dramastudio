import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface GenerationJob {
  id: string;
  capability: string;
  status: string;
  provider?: string;
  model?: string;
  cost?: number;
  created_at: string;
}

export const animationKeys = {
  jobs: (projectId: string) => ["generation-jobs", projectId] as const,
};

export function useGenerationJobs(projectId: string) {
  return useQuery({
    queryKey: animationKeys.jobs(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: GenerationJob[] }>(
        `/v1/projects/${projectId}/ai/jobs`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
