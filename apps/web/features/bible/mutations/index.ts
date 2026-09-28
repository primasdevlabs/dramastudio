import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { SeriesBible } from "@/lib/api/types";
import { bibleKeys } from "../queries";

interface SaveBibleInput {
  projectId: string;
  premise: string;
  genre: string;
  themes: string[];
  tone?: string;
  world_rules: string[];
  narrative_rules: string[];
  visual_direction?: string;
  dialogue_style?: string;
  story_constraints?: string[];
}

export function useSaveBible() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, ...data }: SaveBibleInput) =>
      api.post<SeriesBible>(`/v1/projects/${projectId}/bible`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: bibleKeys.detail(variables.projectId),
      });
    },
  });
}
