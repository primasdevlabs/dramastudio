import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { LeadDirectorDecision, ApprovalRequest } from "@/lib/api/types";
import { productionKeys } from "../queries";

export function useTriggerDirector() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      episodeId,
    }: {
      projectId: string;
      episodeId: string;
    }) =>
      api.post<LeadDirectorDecision>("/v1/agents/lead-director/step", {
        project_id: projectId,
        episode_id: episodeId,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.decisions(variables.projectId),
      });
    },
  });
}

export function useSubmitApproval() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: {
      project_id: string;
      episode_id: string;
      stage: string;
      target_id: string;
      decision: "APPROVE" | "REJECT" | "REQUEST_REVISION" | "REGENERATE" | "PAUSE" | "STOP" | "OVERRIDE";
      notes: string;
      decided_by: string;
    }) => api.post<ApprovalRequest>("/v1/production/approval", input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.approvals(variables.project_id),
      });
    },
  });
}

export function usePauseProduction() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (projectId: string) =>
      api.post(`/v1/production/${projectId}/pause`, {}),
    onSuccess: (_, projectId) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.decisions(projectId),
      });
    },
  });
}

export function useResumeProduction() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (projectId: string) =>
      api.post(`/v1/production/${projectId}/resume`, {}),
    onSuccess: (_, projectId) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.decisions(projectId),
      });
    },
  });
}
