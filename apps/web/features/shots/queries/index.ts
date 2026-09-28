import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Asset } from "@/lib/api/types";

export const shotKeys = {
  all: (projectId: string) => ["shots", projectId] as const,
  detail: (shotId: string) => ["shots", "detail", shotId] as const,
};

export function useShots(projectId: string) {
  return useQuery({
    queryKey: shotKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ assets: Asset[] }>(
        `/v1/media/assets?project_id=${projectId}&type=video`
      );
      return res.assets ?? [];
    },
    enabled: Boolean(projectId),
  });
}
