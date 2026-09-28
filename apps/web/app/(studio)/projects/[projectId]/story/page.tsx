"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Layers,
  Network,
  Clock,
  BookOpen,
  ChevronRight,
} from "lucide-react";
import {
  Paper,
  Group,
  Stack,
  Title,
  Text,
  Badge,
  ThemeIcon,
  Tabs,
  SimpleGrid,
} from "@mantine/core";
import Link from "next/link";
import { api } from "@/lib/api/client";
import type { SeriesBible, Season } from "@/lib/api/types";

export default function StoryWorkspacePage({
  params,
}: {
  params: { projectId: string };
}) {
  const { projectId } = params;

  const { data: bible } = useQuery({
    queryKey: ["bible", projectId],
    queryFn: () =>
      api.get<SeriesBible>(`/v1/projects/${projectId}/bible`).catch(() => null),
  });

  const { data: seasons } = useQuery({
    queryKey: ["seasons", projectId],
    queryFn: async () => {
      const res = await api.get<{ seasons: Season[] }>(
        `/v1/projects/${projectId}/seasons`
      );
      return res.seasons ?? [];
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="terracotta" variant="light" size={44} radius="md">
              <Layers size={22} />
            </ThemeIcon>
            <div>
              <Title order={3} c="white">
                Story Workspace
              </Title>
              <Text size="xs" c="dimmed">
                Narrative hierarchy, story graph, and chronological timeline
              </Text>
            </div>
          </Group>
          <Group gap="sm">
            <Badge color="terracotta" variant="light" size="sm" className="font-mono">
              {seasons?.length ?? 0} Seasons
            </Badge>
          </Group>
        </Group>
      </Paper>

      <Tabs defaultValue="overview" variant="outline">
        <Tabs.List className="border-b border-studio-border">
          <Tabs.Tab value="overview" leftSection={<Layers size={14} />}>
            Overview
          </Tabs.Tab>
          <Tabs.Tab value="graph" leftSection={<Network size={14} />}>
            Story Graph
          </Tabs.Tab>
          <Tabs.Tab value="timeline" leftSection={<Clock size={14} />}>
            Timeline
          </Tabs.Tab>
        </Tabs.List>

        {/* Overview Tab */}
        <Tabs.Panel value="overview" pt="lg">
          <Stack gap="lg">
            {/* Bible Snapshot */}
            <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
              <Group justify="space-between" mb="md">
                <Group gap="xs">
                  <BookOpen size={16} className="text-amber-400" />
                  <Text fw={700} size="sm" c="white">
                    Series Bible Foundation
                  </Text>
                  {bible && (
                    <Badge color="terracotta" variant="light" size="xs" className="font-mono">
                      v{bible.version}
                    </Badge>
                  )}
                </Group>
                <Link
                  href={`/projects/${projectId}/bible`}
                  className="text-xs font-semibold text-studio-accent hover:underline"
                >
                  Edit Bible
                </Link>
              </Group>
              {bible ? (
                <SimpleGrid cols={{ base: 1, md: 2 }} spacing="md">
                  <Paper p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
                    <Text size="xs" c="dimmed" fw={600}>
                      Premise
                    </Text>
                    <Text size="xs" c="white" mt={4}>
                      {bible.premise}
                    </Text>
                  </Paper>
                  <Paper p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
                    <Text size="xs" c="dimmed" fw={600}>
                      Themes
                    </Text>
                    <Group gap={4} mt={4}>
                      {bible.themes?.map((t) => (
                        <Badge
                          key={t}
                          color="terracotta"
                          variant="outline"
                          size="xs"
                        >
                          {t}
                        </Badge>
                      ))}
                    </Group>
                  </Paper>
                </SimpleGrid>
              ) : (
                <Text size="xs" c="dimmed" fs="italic">
                  No Series Bible created yet.
                </Text>
              )}
            </Paper>

            {/* Narrative Hierarchy */}
            <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
              <Text fw={700} size="sm" c="white" mb="md">
                Narrative Hierarchy
              </Text>
              <Stack gap="xs">
                {seasons && seasons.length > 0 ? (
                  seasons.map((season) => (
                    <Paper
                      key={season.id}
                      p="md"
                      radius="sm"
                      className="bg-studio-panel border border-studio-border hover:border-studio-accent/40 transition-colors"
                    >
                      <Group justify="space-between">
                        <Group gap="sm">
                          <Badge color="terracotta" variant="filled" size="sm">
                            S{season.number}
                          </Badge>
                          <div>
                            <Text fw={700} size="sm" c="white">
                              {season.title || `Season ${season.number}`}
                            </Text>
                            <Text size="xs" c="dimmed">
                              {season.episode_ids?.length ?? 0} episodes
                            </Text>
                          </div>
                        </Group>
                        <Link
                          href={`/projects/${projectId}/seasons/${season.id}`}
                          className="text-studio-accent"
                        >
                          <ChevronRight size={16} />
                        </Link>
                      </Group>
                    </Paper>
                  ))
                ) : (
                  <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">
                    No seasons planned yet. Create a season to begin story
                    development.
                  </Text>
                )}
              </Stack>
            </Paper>
          </Stack>
        </Tabs.Panel>

        {/* Story Graph Tab */}
        <Tabs.Panel value="graph" pt="lg">
          <Paper
            p="xl"
            radius="md"
            withBorder
            className="bg-studio-card border-studio-border min-h-[400px] flex items-center justify-center"
          >
            <Stack align="center" gap="md">
              <ThemeIcon
                variant="light"
                color="terracotta"
                size={48}
                radius="md"
              >
                <Network size={24} />
              </ThemeIcon>
              <div className="text-center">
                <Text fw={700} size="md" c="white">
                  Story Dependency Graph
                </Text>
                <Text size="xs" c="dimmed" mt={4} maw={400}>
                  Visualizes narrative dependencies — events, reveals,
                  conflicts, decisions, and their causal relationships.
                </Text>
              </div>
              <Badge color="terracotta" variant="light" size="sm">
                Visual Map
              </Badge>
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Timeline Tab */}
        <Tabs.Panel value="timeline" pt="lg">
          <Paper
            p="xl"
            radius="md"
            withBorder
            className="bg-studio-card border-studio-border"
          >
            <Text fw={700} size="sm" c="white" mb="lg">
              Story Chronology
            </Text>
            <Stack gap="md">
              {seasons && seasons.length > 0 ? (
                seasons.map((season) => (
                  <div key={season.id}>
                    <Group gap="xs" mb="sm">
                      <div className="w-2.5 h-2.5 rounded-full bg-studio-accent" />
                      <Text fw={700} size="xs" c="white" className="font-mono">
                        Season {season.number}: {season.title}
                      </Text>
                    </Group>
                    <div className="ml-5 pl-4 border-l-2 border-studio-border space-y-2">
                      {(season.episode_ids ?? []).map((epId, idx) => (
                        <Paper
                          key={epId}
                          p="xs"
                          radius="sm"
                          className="bg-studio-panel border border-studio-border"
                        >
                          <Text size="xs" c="dimmed" className="font-mono">
                            Episode {idx + 1}
                          </Text>
                        </Paper>
                      ))}
                      {(!season.episode_ids || season.episode_ids.length === 0) && (
                        <Text size="xs" c="dimmed" fs="italic">
                          No episodes planned
                        </Text>
                      )}
                    </div>
                  </div>
                ))
              ) : (
                <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">
                  Timeline populates as seasons and episodes are created.
                </Text>
              )}
            </Stack>
          </Paper>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
