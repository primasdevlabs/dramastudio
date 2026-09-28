import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Location } from "@/lib/api/types";
import { worldKeys } from "../queries";

interface CreateLocationInput {
  project_id: string;
  name: string;
  description: string;
  type: string;
}

export function useCreateLocation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateLocationInput) =>
      api.post<Location>("/v1/world/locations", input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: worldKeys.locations(variables.project_id),
      });
    },
  });
}
