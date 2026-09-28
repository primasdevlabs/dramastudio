import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Character } from "@/lib/api/types";

export const characterKeys = {
  all: (projectId: string) => ["characters", projectId] as const,
  detail: (projectId: string, characterId: string) =>
    ["characters", "detail", projectId, characterId] as const,
};

export function useCharacters(projectId: string) {
  return useQuery({
    queryKey: characterKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: Character[] }>(
        `/v1/projects/${projectId}/characters`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useCharacter(projectId: string, characterId: string) {
  return useQuery({
    queryKey: characterKeys.detail(projectId, characterId),
    queryFn: () =>
      api.get<Character>(
        `/v1/projects/${projectId}/characters/${characterId}`
      ),
    enabled: Boolean(projectId && characterId),
  });
}
