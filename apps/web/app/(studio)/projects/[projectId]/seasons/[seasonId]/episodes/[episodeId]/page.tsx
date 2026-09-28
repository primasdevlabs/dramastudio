"use client";

import { useQuery } from "@tanstack/react-query";
import {
  FileText,
  MessageSquare,
  MapPin,
  Image,
  Video,
  FolderOpen,
  Mic,
  Film,
  ShieldCheck,
  CheckSquare,
  Check,
  Circle,
  Play,
  Bot,
} from "lucide-react";
import {
  Paper,
  Group,
  Stack,
  Text,
  Badge,
  Button,
  Tabs,
  ThemeIcon,
} from "@mantine/core";
import { useParams } from "next/navigation";
import { api } from "@/lib/api/client";
import type { Episode, Scene } from "@/lib/api/types";

const PHASE_STATUS: Record<string, { done: boolean; label: string }> = {
  Script: { done: false, label: "Script" },
  Dialogue: { done: false, label: "Dialogue" },
  Storyboard: { done: false, label: "Storyboard" },
  Assets: { done: false, label: "Assets" },
  Animation: { done: false, label: "Animation" },
  Audio: { done: false, label: "Audio" },
  Assembly: { done: false, label: "Assembly" },
  QA: { done: false, label: "QA" },
};

export default function EpisodeWorkspacePage({
  params,
}: {
  params?: { projectId?: string; seasonId?: string; episodeId?: string };
} = {}) {
  const routeParams = useParams<{ projectId: string; seasonId: string; episodeId: string }>();
  const projectId = params?.projectId ?? routeParams?.projectId ?? "";
  const seasonId = params?.seasonId ?? routeParams?.seasonId ?? "";
  const episodeId = params?.episodeId ?? routeParams?.episodeId ?? "";

  const { data: episode } = useQuery({
    queryKey: ["episodes", "detail", episodeId],
    queryFn: () => api.get<Episode>(`/v1/episodes/${episodeId}`).catch(() => null),
  });

  const { data: scenes } = useQuery({
    queryKey: ["scenes", episodeId],
    queryFn: async () => {
      const res = await api.get<{ scenes: Scene[] }>(
        `/v1/episodes/${episodeId}/scenes`
      );
      return res.scenes ?? [];
    },
  });

  return (
    <Stack gap="md" className="max-w-7xl mx-auto">
      {/* Episode Header */}
      <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="center">
          <div>
            <Group gap="xs" mb={4}>
              <Badge color="terracotta" variant="filled" size="sm" className="font-mono">
                S{episode?.season_id ? "01" : "??"} E
                {String(episode?.number ?? 0).padStart(2, "0")}
              </Badge>
              <Badge
                color={episode?.status === "COMPLETED" ? "emerald" : "amber"}
                variant="light"
                size="xs"
                className="font-mono"
              >
                {episode?.status ?? "PLANNED"}
              </Badge>
            </Group>
            <Text fw={800} size="xl" c="white" className="tracking-tight">
              {episode?.title ?? "Loading Episode..."}
            </Text>
            <Text size="xs" c="dimmed" mt={2}>
              {episode?.summary ?? "Episode production workspace"}
            </Text>
          </div>
          <Group gap="sm">
            <Button
              variant="outline"
              color="terracotta"
              size="xs"
              leftSection={<Bot size={14} />}
            >
              Trigger Director
            </Button>
            <Button
              variant="filled"
              color="terracotta"
              size="xs"
              leftSection={<Play size={14} />}
            >
              Render Episode
            </Button>
          </Group>
        </Group>

        {/* Production Status Bar */}
        <Group gap="xs" mt="md" className="border-t border-studio-border pt-3">
          {Object.entries(PHASE_STATUS).map(([key, phase]) => (
            <Badge
              key={key}
              variant={phase.done ? "filled" : "outline"}
              color={phase.done ? "emerald" : "gray"}
              size="xs"
              className="font-mono"
            >
              {phase.done ? (
                <Check size={10} style={{ verticalAlign: "middle" }} />
              ) : (
                <Circle size={8} style={{ verticalAlign: "middle" }} />
              )}{" "}
              {phase.label}
            </Badge>
          ))}
        </Group>
      </Paper>

      {/* Tabbed Workspace */}
      <Tabs defaultValue="script" variant="outline">
        <Tabs.List className="border-b border-studio-border">
          <Tabs.Tab value="script" leftSection={<FileText size={14} />}>
            Script
          </Tabs.Tab>
          <Tabs.Tab value="dialogue" leftSection={<MessageSquare size={14} />}>
            Dialogue
          </Tabs.Tab>
          <Tabs.Tab value="scenes" leftSection={<MapPin size={14} />}>
            Scenes
          </Tabs.Tab>
          <Tabs.Tab value="storyboard" leftSection={<Image size={14} />}>
            Storyboard
          </Tabs.Tab>
          <Tabs.Tab value="shots" leftSection={<Video size={14} />}>
            Shots
          </Tabs.Tab>
          <Tabs.Tab value="assets" leftSection={<FolderOpen size={14} />}>
            Assets
          </Tabs.Tab>
          <Tabs.Tab value="audio" leftSection={<Mic size={14} />}>
            Audio
          </Tabs.Tab>
          <Tabs.Tab value="assembly" leftSection={<Film size={14} />}>
            Assembly
          </Tabs.Tab>
          <Tabs.Tab value="continuity" leftSection={<ShieldCheck size={14} />}>
            Continuity
          </Tabs.Tab>
          <Tabs.Tab value="qa" leftSection={<CheckSquare size={14} />}>
            QA
          </Tabs.Tab>
        </Tabs.List>

        {/* Script Panel */}
        <Tabs.Panel value="script" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Episode Script
              </Text>
              <Group gap="sm">
                <Button variant="outline" color="terracotta" size="xs">
                  Generate Script
                </Button>
                <Button variant="light" color="emerald" size="xs">
                  Approve
                </Button>
              </Group>
            </Group>
            <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border min-h-[300px]">
              <Text size="xs" c="dimmed" fs="italic">
                Script draft is empty. Click "Draft Screenplay" to write the episode screenplay from the Series Bible and season arc.
              </Text>
            </Paper>
          </Paper>
        </Tabs.Panel>

        {/* Dialogue Panel */}
        <Tabs.Panel value="dialogue" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Dialogue Lines
              </Text>
              <Button variant="outline" color="terracotta" size="xs">
                Generate Dialogue
              </Button>
            </Group>
            <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border min-h-[200px]">
              <Text size="xs" c="dimmed" fs="italic">
                Dialogue is a separately generated artifact. Each line includes
                character, emotion, intent, and delivery direction.
              </Text>
            </Paper>
          </Paper>
        </Tabs.Panel>

        {/* Scenes Panel */}
        <Tabs.Panel value="scenes" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Text fw={700} size="sm" c="white" mb="md">
              Episode Scenes
            </Text>
            <Stack gap="sm">
              {scenes && scenes.length > 0 ? (
                scenes.map((scene) => (
                  <Paper
                    key={scene.id}
                    p="md"
                    radius="sm"
                    className="bg-studio-panel border border-studio-border"
                  >
                    <Group justify="space-between">
                      <Group gap="sm">
                        <Badge color="terracotta" variant="filled" size="sm">
                          Scene {scene.number}
                        </Badge>
                        <Text fw={700} size="xs" c="white">
                          {scene.title}
                        </Text>
                      </Group>
                      <Group gap="xs">
                        <Badge color="gray" variant="outline" size="xs">
                          {scene.time_of_day}
                        </Badge>
                        <Badge color="gray" variant="outline" size="xs">
                          {scene.character_ids?.length ?? 0} characters
                        </Badge>
                      </Group>
                    </Group>
                    <Text size="xs" c="dimmed" mt={4}>
                      {scene.description}
                    </Text>
                  </Paper>
                ))
              ) : (
                <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">
                  No scenes defined for this episode yet.
                </Text>
              )}
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Storyboard Panel */}
        <Tabs.Panel value="storyboard" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Storyboard
              </Text>
              <Group gap="sm">
                <Button variant="outline" color="terracotta" size="xs">
                  Generate Storyboard
                </Button>
                <Button variant="light" color="emerald" size="xs">
                  Approve
                </Button>
              </Group>
            </Group>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[300px] flex items-center justify-center"
            >
              <Stack align="center" gap="sm">
                <ThemeIcon variant="light" color="terracotta" size={40} radius="md">
                  <Image size={20} />
                </ThemeIcon>
                <Text size="xs" c="dimmed">
                  Storyboard panels appear here after generation
                </Text>
              </Stack>
            </Paper>
          </Paper>
        </Tabs.Panel>

        {/* Shots Panel */}
        <Tabs.Panel value="shots" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Video Shots
              </Text>
              <Button variant="outline" color="terracotta" size="xs">
                Generate Video
              </Button>
            </Group>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[300px] flex items-center justify-center"
            >
              <Stack align="center" gap="sm">
                <ThemeIcon variant="light" color="terracotta" size={40} radius="md">
                  <Video size={20} />
                </ThemeIcon>
                <Text size="xs" c="dimmed">
                  Generated video shots with 9:16 vertical previews
                </Text>
              </Stack>
            </Paper>
          </Paper>
        </Tabs.Panel>

        {/* Assets Panel */}
        <Tabs.Panel value="assets" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Text fw={700} size="sm" c="white" mb="md">
              Episode Assets
            </Text>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[200px] flex items-center justify-center"
            >
              <Text size="xs" c="dimmed">
                All generated media assets for this episode — images, video,
                audio, renders.
              </Text>
            </Paper>
          </Paper>
        </Tabs.Panel>

        {/* Audio Panel */}
        <Tabs.Panel value="audio" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Audio Tracks
              </Text>
              <Button variant="outline" color="terracotta" size="xs">
                Generate Audio
              </Button>
            </Group>
            <Stack gap="xs">
              {["Dialogue", "Music / Score", "Sound Effects"].map((track) => (
                <Paper
                  key={track}
                  p="sm"
                  radius="sm"
                  className="bg-studio-panel border border-studio-border"
                >
                  <Group justify="space-between">
                    <Text size="xs" fw={600} c="white">
                      {track}
                    </Text>
                    <Badge color="gray" variant="outline" size="xs">
                      Not generated
                    </Badge>
                  </Group>
                </Paper>
              ))}
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Assembly Panel */}
        <Tabs.Panel value="assembly" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Assembly Timeline
              </Text>
              <Button variant="filled" color="terracotta" size="xs">
                Render Episode
              </Button>
            </Group>
            <Stack gap="sm">
              {[
                { label: "Video Track", color: "studio-accent" },
                { label: "Dialogue Track", color: "amber-400" },
                { label: "Music Track", color: "emerald-400" },
                { label: "SFX Track", color: "studio-muted" },
              ].map((track) => (
                <div key={track.label} className="space-y-1">
                  <Text size="xs" c="dimmed" fw={600}>
                    {track.label}
                  </Text>
                  <div
                    className={`h-10 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto`}
                  >
                    <div className="h-full bg-studio-accent/10 border border-studio-accent/20 rounded px-3 flex items-center text-xs text-studio-muted font-mono min-w-[160px]">
                      No clips
                    </div>
                  </div>
                </div>
              ))}
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Continuity Panel */}
        <Tabs.Panel value="continuity" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Text fw={700} size="sm" c="white" mb="md">
              Episode Continuity
            </Text>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[200px] flex items-center justify-center"
            >
              <Stack align="center" gap="sm">
                <ThemeIcon variant="light" color="emerald" size={40} radius="md">
                  <ShieldCheck size={20} />
                </ThemeIcon>
                <Text size="xs" c="dimmed">
                  Continuity checks run automatically. Issues appear here when
                  detected.
                </Text>
              </Stack>
            </Paper>
          </Paper>
        </Tabs.Panel>

        {/* QA Panel */}
        <Tabs.Panel value="qa" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Text fw={700} size="sm" c="white" mb="md">
              Quality Assurance
            </Text>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[200px] flex items-center justify-center"
            >
              <Stack align="center" gap="sm">
                <ThemeIcon variant="light" color="emerald" size={40} radius="md">
                  <CheckSquare size={20} />
                </ThemeIcon>
                <Text size="xs" c="dimmed">
                  Final quality review before episode is marked ready for
                  rendering and publishing.
                </Text>
              </Stack>
            </Paper>
          </Paper>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
