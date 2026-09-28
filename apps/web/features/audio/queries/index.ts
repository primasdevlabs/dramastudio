import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface AudioTrack {
  id: string;
  episode_id: string;
  type: "dialogue" | "music" | "sfx";
  character_id?: string;
  url: string;
  duration_seconds: number;
  status: "PENDING" | "GENERATED" | "APPROVED";
}

export const audioKeys = {
  all: (episodeId: string) => ["audio", episodeId] as const,
};

export function useAudioTracks(episodeId: string) {
  return useQuery({
    queryKey: audioKeys.all(episodeId),
    queryFn: async () => {
      const res = await api.get<{ tracks: AudioTrack[] }>(
        `/v1/episodes/${episodeId}/audio`
      );
      return res.tracks ?? [];
    },
    enabled: Boolean(episodeId),
  });
}
