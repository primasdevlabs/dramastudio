import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { continuityKeys } from "../queries";

export function useResolveIssue() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      issueId,
      resolution,
    }: {
      projectId: string;
      issueId: string;
      resolution: string;
    }) =>
      api.post(
        `/v1/projects/${projectId}/continuity/issues/${issueId}/resolve`,
        { resolution }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: continuityKeys.all(variables.projectId),
      });
    },
  });
}
