import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface DialogueLine {
  id: string;
  character: string;
  line: string;
  emotion: string;
  intent: string;
  delivery: string;
  scene_number: number;
}

export const dialogueKeys = {
  all: (episodeId: string) => ["dialogue", episodeId] as const,
};

export function useDialogue(episodeId: string) {
  return useQuery({
    queryKey: dialogueKeys.all(episodeId),
    queryFn: async () => {
      const res = await api.get<{ lines: DialogueLine[] }>(
        `/v1/episodes/${episodeId}/dialogue`
      );
      return res.lines ?? [];
    },
    enabled: Boolean(episodeId),
  });
}
