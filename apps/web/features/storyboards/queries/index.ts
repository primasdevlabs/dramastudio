import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface Storyboard {
  id: string;
  episode_id: string;
  scene_id: string;
  version: number;
  shots: StoryboardShot[];
  status: "DRAFT" | "APPROVED" | "REJECTED";
}

export interface StoryboardShot {
  id: string;
  number: number;
  shot_type: string;
  camera: string;
  character_ids: string[];
  action: string;
  location_id: string;
  duration_seconds: number;
  thumbnail_url?: string;
  status: "PENDING" | "GENERATED" | "APPROVED" | "REJECTED";
}

export const storyboardKeys = {
  detail: (episodeId: string) => ["storyboard", episodeId] as const,
  shots: (storyboardId: string) => ["storyboard-shots", storyboardId] as const,
};

export function useStoryboard(episodeId: string) {
  return useQuery({
    queryKey: storyboardKeys.detail(episodeId),
    queryFn: () =>
      api.get<Storyboard>(`/v1/episodes/${episodeId}/storyboard`).catch(() => null),
    enabled: Boolean(episodeId),
  });
}
