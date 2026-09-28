import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Project } from "@/lib/api/types";
import { projectKeys } from "../queries";

interface CreateProjectInput {
  name: string;
  description: string;
  genre: string;
  language: string;
  mode: "monitored" | "autonomous";
}

export function useCreateProject() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateProjectInput) =>
      api.post<Project>("/v1/projects", input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.all });
    },
  });
}

