import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { StoryFact } from "@/lib/api/types";
import { canonKeys } from "../queries";

interface CreateFactInput {
  project_id: string;
  subject: string;
  predicate: string;
  object: string;
  introduced: string;
  valid_from: string;
}

export function useCreateFact() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ project_id, ...data }: CreateFactInput) =>
      api.post<StoryFact>(`/v1/projects/${project_id}/canon/facts`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: canonKeys.all(variables.project_id),
      });
    },
  });
}
