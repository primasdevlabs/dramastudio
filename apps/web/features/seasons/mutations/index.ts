import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Season } from "@/lib/api/types";
import { seasonKeys } from "../queries";

interface CreateSeasonInput {
  project_id: string;
  number: number;
  title: string;
  summary: string;
}

export function useCreateSeason() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ project_id, ...data }: CreateSeasonInput) =>
      api.post<Season>(`/v1/projects/${project_id}/seasons`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: seasonKeys.all(variables.project_id),
      });
    },
  });
}
