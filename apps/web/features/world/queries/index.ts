import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { Location } from "@/lib/api/types";

export const worldKeys = {
  locations: (projectId: string) => ["locations", projectId] as const,
};

export function useLocations(projectId: string) {
  return useQuery({
    queryKey: worldKeys.locations(projectId),
    queryFn: async () => {
      const res = await api.get<{ items: Location[] }>(
        `/v1/projects/${projectId}/world/locations`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });
}
