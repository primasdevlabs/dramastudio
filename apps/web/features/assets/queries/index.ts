import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Asset } from "@/lib/api/types";

export const assetKeys = {
  all: (projectId: string) => ["assets", projectId] as const,
  byType: (projectId: string, type: string) =>
    ["assets", projectId, type] as const,
};

export function useAssets(projectId: string) {
  return useQuery({
    queryKey: assetKeys.all(projectId),
    queryFn: async () => {
      const res = await api.get<{ assets: Asset[] }>(
        `/v1/media/assets?project_id=${projectId}`
      );
      return res.assets ?? [];
    },
    enabled: Boolean(projectId),
  });
}
