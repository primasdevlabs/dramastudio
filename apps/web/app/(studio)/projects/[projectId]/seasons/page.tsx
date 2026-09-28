"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Film,
  Plus,
  ChevronRight,
  Calendar,
} from "lucide-react";
import {
  Paper,
  Group,
  Stack,
  Title,
  Text,
  Badge,
  Button,
  Modal,
  TextInput,
  Textarea,
  ThemeIcon,
  Card,
  SimpleGrid,
  NumberInput,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import Link from "next/link";
import { api } from "@/lib/api/client";
import type { Season } from "@/lib/api/types";

export default function SeasonsPage({
  params,
}: {
  params: { projectId: string };
}) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [opened, setOpened] = useState(false);
  const [title, setTitle] = useState("");
  const [summary, setSummary] = useState("");
  const [number, setNumber] = useState<number>(1);

  const { data: seasons, isLoading } = useQuery({
    queryKey: ["seasons", projectId],
    queryFn: async () => {
      const res = await api.get<{ seasons: Season[] }>(
        `/v1/projects/${projectId}/seasons`
      );
      return res.seasons ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: () =>
      api.post<Season>(`/v1/projects/${projectId}/seasons`, {
        number,
        title,
        summary,
      }),
    onSuccess: (s) => {
      queryClient.invalidateQueries({ queryKey: ["seasons", projectId] });
      setOpened(false);
      setTitle("");
      setSummary("");
      notifications.show({
        title: "Season Created",
        message: `Added Season ${s.number || number}: ${s.title || title}`,
        color: "terracotta",
      });
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="terracotta" variant="light" size={44} radius="md">
              <Calendar size={22} />
            </ThemeIcon>
            <div>
              <Title order={3} c="white">
                Seasons
              </Title>
              <Text size="xs" c="dimmed">
                Plan and manage serialized season arcs and episode structure
              </Text>
            </div>
          </Group>
          <Button
            onClick={() => setOpened(true)}
            leftSection={<Plus size={16} />}
            variant="filled"
            color="terracotta"
            size="sm"
            radius="sm"
          >
            Plan New Season
          </Button>
        </Group>
      </Paper>

      {/* Create Season Modal */}
      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title={
          <Text fw={700} size="md" c="white">
            Plan New Season
          </Text>
        }
        centered
        overlayProps={{ backgroundOpacity: 0.7, blur: 4 }}
      >
        <Stack gap="md">
          <NumberInput
            label="Season Number"
            value={number}
            onChange={(val) => setNumber(typeof val === "number" ? val : 1)}
            min={1}
            max={99}
            variant="filled"
          />
          <TextInput
            label="Season Title"
            placeholder="e.g. The Reckoning"
            required
            value={title}
            onChange={(e) => setTitle(e.currentTarget.value)}
            variant="filled"
          />
          <Textarea
            label="Season Summary"
            placeholder="Outline the season arc, key conflicts, and narrative goals..."
            rows={4}
            value={summary}
            onChange={(e) => setSummary(e.currentTarget.value)}
            variant="filled"
          />
          <Group justify="flex-end" gap="sm" mt="md">
            <Button variant="subtle" color="gray" onClick={() => setOpened(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => createMutation.mutate()}
              loading={createMutation.isPending}
              disabled={!title}
              color="terracotta"
            >
              Create Season
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Seasons Grid */}
      <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="lg">
        {isLoading ? (
          [1, 2].map((i) => (
            <Paper key={i} h={180} radius="md" className="bg-studio-card animate-pulse" />
          ))
        ) : seasons && seasons.length > 0 ? (
          seasons.map((season) => (
            <Card
              key={season.id}
              p="lg"
              radius="md"
              withBorder
              className="bg-studio-card border-studio-border hover:border-studio-accent/60 transition-all flex flex-col justify-between"
            >
              <Stack gap="sm">
                <Group justify="space-between">
                  <Badge color="terracotta" variant="filled" size="lg">
                    Season {season.number}
                  </Badge>
                  <Badge color="emerald" variant="dot" size="sm">
                    {season.episode_ids?.length ?? 0} Episodes
                  </Badge>
                </Group>
                <div>
                  <Text fw={700} size="lg" c="white" className="tracking-tight">
                    {season.title || `Season ${season.number}`}
                  </Text>
                  <Text size="xs" c="dimmed" className="line-clamp-2 mt-1">
                    {season.summary || "No season summary specified."}
                  </Text>
                </div>
              </Stack>
              <Group justify="flex-end" className="pt-3 border-t border-studio-border/60 mt-4">
                <Button
                  component={Link}
                  href={`/projects/${projectId}/seasons/${season.id}`}
                  variant="subtle"
                  color="terracotta"
                  size="xs"
                  rightSection={<ChevronRight size={14} />}
                >
                  Open Season
                </Button>
              </Group>
            </Card>
          ))
        ) : (
          <Paper
            p="xl"
            radius="md"
            withBorder
            className="col-span-full bg-studio-card border-studio-border text-center"
          >
            <Stack align="center" gap="md">
              <ThemeIcon variant="light" color="terracotta" size={48} radius="md">
                <Calendar size={24} />
              </ThemeIcon>
              <Text size="xs" c="dimmed">
                No seasons planned yet. Create a season to start structuring your
                story.
              </Text>
            </Stack>
          </Paper>
        )}
      </SimpleGrid>
    </Stack>
  );
}
