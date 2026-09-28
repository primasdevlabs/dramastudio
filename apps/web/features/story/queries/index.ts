import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";

export interface StoryNode {
  id: string;
  type: "event" | "reveal" | "conflict" | "decision" | "relationship_change" | "foreshadowing" | "resolution";
  label: string;
  episode_id: string;
  description: string;
}

export interface StoryEdge {
  id: string;
  source: string;
  target: string;
  relationship: "causes" | "reveals" | "depends_on" | "contradicts" | "resolves" | "foreshadows";
}

export interface StoryGraph {
  nodes: StoryNode[];
  edges: StoryEdge[];
}

export const storyKeys = {
  graph: (projectId: string) => ["story-graph", projectId] as const,
};

export function useStoryGraph(projectId: string) {
  return useQuery({
    queryKey: storyKeys.graph(projectId),
    queryFn: () =>
      api
        .get<StoryGraph>(`/v1/projects/${projectId}/story/graph`)
        .catch(() => ({ nodes: [], edges: [] })),
    enabled: Boolean(projectId),
  });
}
