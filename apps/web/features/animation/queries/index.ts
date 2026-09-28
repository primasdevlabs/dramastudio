import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface GenerationJob {
  id: string;
  shot_id: string;
  type: "video" | "image" | "audio";
  status: "QUEUED" | "GENERATING" | "PROCESSING" | "COMPLETED" | "FAILED";
  provider: string;
  model: string;
  progress_percent?: number;
  created_at: string;
}

export const animationKeys = {
  jobs: (projectId: string) => ["animation-jobs", projectId] as const,
};

export function useGenerationJobs(projectId: string) {
  return useQuery({
    queryKey: animationKeys.jobs(projectId),
    queryFn: async () => {
      const res = await api.get<{ jobs: GenerationJob[] }>(
        `/v1/media/jobs?project_id=${projectId}`
      );
      return res.jobs ?? [];
    },
    enabled: Boolean(projectId),
  });
}
