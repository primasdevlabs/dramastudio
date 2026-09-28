import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface Script {
  id: string;
  episode_id: string;
  version: number;
  content: string;
  scenes: { number: number; beats: string[]; action: string; dialogue: string }[];
  status: "DRAFT" | "APPROVED" | "REJECTED";
  created_at: string;
}

export const scriptKeys = {
  detail: (episodeId: string) => ["script", episodeId] as const,
};

export function useScript(episodeId: string) {
  return useQuery({
    queryKey: scriptKeys.detail(episodeId),
    queryFn: () =>
      api.get<Script>(`/v1/episodes/${episodeId}/script`).catch(() => null),
    enabled: Boolean(episodeId),
  });
}
