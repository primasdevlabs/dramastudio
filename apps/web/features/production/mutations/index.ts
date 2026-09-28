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
      api.post<LeadDirectorDecision>(
        `/v1/projects/${projectId}/agents/director/step`,
        { episode_id: episodeId }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.decisions(variables.projectId),
      });
    },
  });
}

// Approvals are a two-step flow: request the approval, then decide it.
// decided_by is derived from the authenticated principal server-side.
export function useSubmitApproval() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: {
      project_id: string;
      episode_id: string;
      stage: string;
      target_id: string;
      decision: "APPROVE" | "REJECT" | "REQUEST_REVISION" | "REGENERATE" | "PAUSE" | "STOP" | "OVERRIDE";
      notes: string;
    }) => {
      const approval = await api.post<ApprovalRequest>(
        `/v1/projects/${input.project_id}/production/approvals`,
        {
          episode_id: input.episode_id,
          stage: input.stage,
          target_id: input.target_id,
        }
      );
      return api.post<ApprovalRequest>(
        `/v1/projects/${input.project_id}/production/approvals/${approval.id}/decide`,
        { decision: input.decision, notes: input.notes }
      );
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.approvals(variables.project_id),
      });
    },
  });
}

// Pause/resume act on an individual production run, not the project.
export function usePauseRun() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, runId }: { projectId: string; runId: string }) =>
      api.post(`/v1/projects/${projectId}/production/runs/${runId}/pause`, {}),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.runs(variables.projectId),
      });
    },
  });
}

export function useResumeRun() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, runId }: { projectId: string; runId: string }) =>
      api.post(`/v1/projects/${projectId}/production/runs/${runId}/resume`, {}),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: productionKeys.runs(variables.projectId),
      });
    },
  });
}
