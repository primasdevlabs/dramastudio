import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Character } from "@/lib/api/types";

export const characterKeys = {
  all: (projectId: string) => ["characters", projectId] as const,
  detail: (characterId: string) => ["characters", "detail", characterId] as const,
};

export function useCharacters(projectId: string) {
  return useQuery({
    queryKey: characterKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ characters: Character[] }>(
        `/v1/characters?project_id=${projectId}`
      );
      return res.characters ?? [];
    },
    enabled: Boolean(projectId),
  });
}

export function useCharacter(characterId: string) {
  return useQuery({
    queryKey: characterKeys.detail(characterId),
    queryFn: () => api.get<Character>(`/v1/characters/${characterId}`),
    enabled: Boolean(characterId),
  });
}
