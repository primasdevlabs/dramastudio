import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { assetKeys } from "../queries";

export function useApproveAsset() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      assetId,
    }: {
      projectId: string;
      assetId: string;
    }) => api.post(`/v1/media/assets/${assetId}/approve`, {}),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: assetKeys.all(variables.projectId),
      });
    },
  });
}
