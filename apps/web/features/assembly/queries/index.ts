import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface AssemblyTimeline {
  id: string;
  episode_id: string;
  tracks: AssemblyTrack[];
}

export interface AssemblyTrack {
  type: "video" | "dialogue" | "music" | "sfx" | "subtitles";
  clips: { asset_id: string; start_ms: number; end_ms: number; label: string }[];
}

export const assemblyKeys = {
  detail: (episodeId: string) => ["assembly", episodeId] as const,
};

export function useAssembly(episodeId: string) {
  return useQuery({
    queryKey: assemblyKeys.detail(episodeId),
    queryFn: () =>
      api
        .get<AssemblyTimeline>(`/v1/episodes/${episodeId}/assembly`)
        .catch(() => null),
    enabled: Boolean(episodeId),
  });
}
