import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Character } from "@/lib/api/types";
import { characterKeys } from "../queries";

interface CreateCharacterInput {
  project_id: string;
  name: string;
  role: string;
  bio: string;
}

export function useCreateCharacter() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateCharacterInput) =>
      api.post<Character>("/v1/characters", input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: characterKeys.all(variables.project_id),
      });
    },
  });
}
