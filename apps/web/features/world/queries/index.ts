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
      const res = await api.get<{ locations: Location[] }>(
        `/v1/world/locations?project_id=${projectId}`
      );
      return res.locations ?? [];
    },
    enabled: Boolean(projectId),
  });
}
