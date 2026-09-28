"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Film, Plus, ChevronRight, PlayCircle } from "lucide-react";
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
  NumberInput,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import Link from "next/link";
import { useParams } from "next/navigation";
import { api } from "@/lib/api/client";
import type { Season, Episode } from "@/lib/api/types";

const STATUS_COLORS: Record<string, string> = {
  PLANNED: "gray",
  SCRIPTED: "amber",
  PRODUCING: "terracotta",
  VALIDATING: "amber",
  COMPLETED: "emerald",
};

export default function SeasonDetailPage() {
  const { projectId, seasonId } = useParams<{ projectId: string; seasonId: string }>();
  const queryClient = useQueryClient();

  const [opened, setOpened] = useState(false);
  const [epTitle, setEpTitle] = useState("");
  const [epSummary, setEpSummary] = useState("");
  const [epNumber, setEpNumber] = useState<number>(1);

  const { data: season } = useQuery({
    queryKey: ["seasons", "detail", seasonId],
    queryFn: () => api.get<Season>(`/v1/projects/${projectId}/seasons/${seasonId}`),
  });

  const { data: episodes, isLoading } = useQuery({
    queryKey: ["episodes", seasonId],
    queryFn: async () => {
      const res = await api.get<{ items: Episode[] }>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes`
      );
      return res.items ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: () =>
      api.post<Episode>(`/v1/projects/${projectId}/seasons/${seasonId}/episodes`, {
        number: epNumber,
        title: epTitle,
        summary: epSummary,
      }),
    onSuccess: (ep) => {
      queryClient.invalidateQueries({ queryKey: ["episodes", seasonId] });
      setOpened(false);
      setEpTitle("");
      setEpSummary("");
      notifications.show({
        title: "Episode Created",
        message: `Added Episode ${ep.number || epNumber}: ${ep.title || epTitle}`,
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
              <Film size={22} />
            </ThemeIcon>
            <div>
              <Group gap="xs">
                <Badge color="terracotta" variant="filled" size="sm">
                  Season {season?.number}
                </Badge>
                <Title order={3} c="white">
                  {season?.title || "Loading..."}
                </Title>
              </Group>
              <Text size="xs" c="dimmed">
                {season?.summary || "Season arc and episode planning"}
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
            Add Episode
          </Button>
        </Group>
      </Paper>

      {/* Create Episode Modal */}
      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title={
          <Text fw={700} size="md" c="white">
            Add Episode
          </Text>
        }
        centered
        overlayProps={{ backgroundOpacity: 0.7, blur: 4 }}
      >
        <Stack gap="md">
          <NumberInput
            label="Episode Number"
            value={epNumber}
            onChange={(val) => setEpNumber(typeof val === "number" ? val : 1)}
            min={1}
            variant="filled"
          />
          <TextInput
            label="Episode Title"
            placeholder="e.g. Midnight Call"
            required
            value={epTitle}
            onChange={(e) => setEpTitle(e.currentTarget.value)}
            variant="filled"
          />
          <Textarea
            label="Episode Summary"
            placeholder="Brief episode synopsis..."
            rows={3}
            value={epSummary}
            onChange={(e) => setEpSummary(e.currentTarget.value)}
            variant="filled"
          />
          <Group justify="flex-end" gap="sm" mt="md">
            <Button variant="subtle" color="gray" onClick={() => setOpened(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => createMutation.mutate()}
              loading={createMutation.isPending}
              disabled={!epTitle}
              color="terracotta"
            >
              Create Episode
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Episodes List */}
      <Stack gap="md">
        {isLoading ? (
          [1, 2, 3].map((i) => (
            <Paper key={i} h={80} radius="md" className="bg-studio-card animate-pulse" />
          ))
        ) : episodes && episodes.length > 0 ? (
          episodes.map((ep) => (
            <Card
              key={ep.id}
              p="lg"
              radius="md"
              withBorder
              className="bg-studio-card border-studio-border hover:border-studio-accent/60 transition-all"
            >
              <Group justify="space-between" align="center">
                <Group gap="md">
                  <div className="w-12 h-12 rounded bg-studio-panel border border-studio-border flex items-center justify-center">
                    <Text fw={800} size="lg" c="white" className="font-mono">
                      {String(ep.number).padStart(2, "0")}
                    </Text>
                  </div>
                  <div>
                    <Text fw={700} size="sm" c="white">
                      {ep.title || `Episode ${ep.number}`}
                    </Text>
                    <Text size="xs" c="dimmed" className="line-clamp-1">
                      {ep.summary || "No synopsis"}
                    </Text>
                  </div>
                </Group>
                <Group gap="md">
                  <Badge
                    color={STATUS_COLORS[ep.status] ?? "gray"}
                    variant="light"
                    size="sm"
                    className="font-mono"
                  >
                    {ep.status}
                  </Badge>
                  <Badge color="gray" variant="outline" size="xs" className="font-mono">
                    {ep.scene_ids?.length ?? 0} scenes
                  </Badge>
                  <Button
                    component={Link}
                    href={`/projects/${projectId}/seasons/${seasonId}/episodes/${ep.id}`}
                    variant="subtle"
                    color="terracotta"
                    size="xs"
                    rightSection={<ChevronRight size={14} />}
                  >
                    Open
                  </Button>
                </Group>
              </Group>
            </Card>
          ))
        ) : (
          <Paper
            p="xl"
            radius="md"
            withBorder
            className="bg-studio-card border-studio-border text-center"
          >
            <Stack align="center" gap="md">
              <ThemeIcon variant="light" color="terracotta" size={48} radius="md">
                <PlayCircle size={24} />
              </ThemeIcon>
              <Text size="xs" c="dimmed">
                No episodes created yet. Add an episode to begin production.
              </Text>
            </Stack>
          </Paper>
        )}
      </Stack>
    </Stack>
  );
}
