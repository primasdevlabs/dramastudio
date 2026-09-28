"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
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
import type { Episode, Scene, Shot, Asset, EpisodeTimeline, ContinuityIssue } from "@/lib/api/types";

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
    queryFn: () =>
      api
        .get<Episode>(`/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}`)
        .catch(() => null),
    enabled: Boolean(projectId && seasonId && episodeId),
  });

  const { data: scenes } = useQuery({
    queryKey: ["scenes", episodeId],
    queryFn: async () => {
      const res = await api.get<{ items: Scene[] }>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}/scenes`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId && seasonId && episodeId),
  });

  const queryClient = useQueryClient();
  const invalidateEpisode = () =>
    queryClient.invalidateQueries({ queryKey: ["episodes", "detail", episodeId] });

  const { data: shots } = useQuery({
    queryKey: ["shots", projectId, episodeId],
    queryFn: async () => {
      const res = await api.get<{ items: Shot[] }>(
        `/v1/projects/${projectId}/production/shots?episode_id=${episodeId}`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId && episodeId),
  });

  const { data: assets } = useQuery({
    queryKey: ["assets", projectId, "episode", episodeId],
    queryFn: async () => {
      const res = await api.get<{ items: Asset[] }>(
        `/v1/projects/${projectId}/media/assets`
      );
      return (res.items ?? []).filter((a) => a.episode_id === episodeId);
    },
    enabled: Boolean(projectId && episodeId),
  });

  const { data: timeline } = useQuery({
    queryKey: ["timeline", projectId, episodeId],
    queryFn: () =>
      api
        .get<EpisodeTimeline>(
          `/v1/projects/${projectId}/postproduction/timelines?episode_id=${episodeId}`
        )
        .catch(() => null),
    enabled: Boolean(projectId && episodeId),
  });

  const { data: continuityIssues } = useQuery({
    queryKey: ["continuity-issues", projectId, episodeId],
    queryFn: async () => {
      const res = await api.get<{ items: ContinuityIssue[] }>(
        `/v1/projects/${projectId}/continuity/issues`
      );
      return (res.items ?? []).filter((i) => i.episode_id === episodeId);
    },
    enabled: Boolean(projectId && episodeId),
  });

  const aiGenerate = useMutation({
    mutationFn: (capability: string) =>
      api.post(`/v1/projects/${projectId}/ai/generate`, {
        capability,
        input: { episode_id: episodeId },
      }),
    onSuccess: invalidateEpisode,
  });

  const generateAsset = useMutation({
    mutationFn: (vars: { capability: string; type: string; prompt: string }) =>
      api.post(`/v1/projects/${projectId}/media/assets`, {
        ...vars,
        episode_id: episodeId,
      }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["assets", projectId] }),
  });

  const approveScript = useMutation({
    mutationFn: () =>
      api.patch<Episode>(
        `/v1/projects/${projectId}/seasons/${seasonId}/episodes/${episodeId}`,
        { status: "SCRIPTED" }
      ),
    onSuccess: invalidateEpisode,
  });

  const approveShot = useMutation({
    mutationFn: (shotId: string) =>
      api.post(`/v1/projects/${projectId}/production/shots/${shotId}/approve`, {}),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["shots", projectId, episodeId] }),
  });

  const directorStep = useMutation({
    mutationFn: () =>
      api.post(`/v1/projects/${projectId}/agents/director/step`, {
        episode_id: episodeId,
      }),
  });

  const renderEpisode = useMutation({
    mutationFn: () =>
      api.post(`/v1/projects/${projectId}/postproduction/renders`, {
        episode_id: episodeId,
        format: "mp4",
        resolution: "1080x1920",
      }),
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
              loading={directorStep.isPending}
              onClick={() => directorStep.mutate()}
            >
              Trigger Director
            </Button>
            <Button
              variant="filled"
              color="terracotta"
              size="xs"
              leftSection={<Play size={14} />}
              loading={renderEpisode.isPending}
              onClick={() => renderEpisode.mutate()}
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
                <Button
                  variant="outline"
                  color="terracotta"
                  size="xs"
                  loading={aiGenerate.isPending}
                  onClick={() => aiGenerate.mutate("script_writing")}
                >
                  Generate Script
                </Button>
                <Button
                  variant="light"
                  color="emerald"
                  size="xs"
                  disabled={!episode?.script}
                  loading={approveScript.isPending}
                  onClick={() => approveScript.mutate()}
                >
                  Approve
                </Button>
              </Group>
            </Group>
            <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border min-h-[300px]">
              {episode?.script ? (
                <Text size="xs" c="white" style={{ whiteSpace: "pre-wrap" }} className="font-mono">
                  {episode.script}
                </Text>
              ) : (
                <Text size="xs" c="dimmed" fs="italic">
                  Script draft is empty. Click "Generate Script" to write the episode screenplay from the Series Bible and season arc.
                </Text>
              )}
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
              <Button
                variant="outline"
                color="terracotta"
                size="xs"
                loading={aiGenerate.isPending}
                onClick={() => aiGenerate.mutate("dialogue_writing")}
              >
                Generate Dialogue
              </Button>
            </Group>
            <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border min-h-[200px]">
              <Text size="xs" c="dimmed" fs="italic">
                Dialogue is generated into the episode script. Each line includes
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
                <Button
                  variant="outline"
                  color="terracotta"
                  size="xs"
                  loading={aiGenerate.isPending}
                  onClick={() => aiGenerate.mutate("storyboard")}
                >
                  Generate Storyboard
                </Button>
              </Group>
            </Group>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[300px]"
            >
              {(shots ?? []).length > 0 ? (
                <Stack gap="xs">
                  {(shots ?? []).map((shot) => (
                    <Paper key={shot.id} p="sm" radius="sm" className="bg-studio-card border border-studio-border">
                      <Group justify="space-between">
                        <Group gap="sm">
                          <Badge color="terracotta" variant="outline" size="xs" className="font-mono">
                            #{shot.seq}
                          </Badge>
                          <Text size="xs" c="white">{shot.description}</Text>
                        </Group>
                        <Group gap="xs">
                          <Badge color={shot.status === "approved" ? "emerald" : "gray"} variant="light" size="xs">
                            {shot.status}
                          </Badge>
                          {shot.status !== "approved" && (
                            <Button
                              variant="subtle"
                              color="emerald"
                              size="compact-xs"
                              loading={approveShot.isPending}
                              onClick={() => approveShot.mutate(shot.id)}
                            >
                              Approve
                            </Button>
                          )}
                        </Group>
                      </Group>
                    </Paper>
                  ))}
                </Stack>
              ) : (
                <Stack align="center" gap="sm" py="xl">
                  <ThemeIcon variant="light" color="terracotta" size={40} radius="md">
                    <Image size={20} />
                  </ThemeIcon>
                  <Text size="xs" c="dimmed">
                    Storyboard shots appear here after generation
                  </Text>
                </Stack>
              )}
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
              <Button
                variant="outline"
                color="terracotta"
                size="xs"
                loading={generateAsset.isPending}
                onClick={() =>
                  generateAsset.mutate({
                    capability: "video_generation",
                    type: "video",
                    prompt: `Episode ${episode?.number ?? ""}: ${episode?.title ?? ""}`,
                  })
                }
              >
                Generate Video
              </Button>
            </Group>
            <Paper
              p="xl"
              radius="sm"
              className="bg-studio-panel border border-studio-border min-h-[300px]"
            >
              {(assets ?? []).filter((a) => a.type === "video").length > 0 ? (
                <Stack gap="xs">
                  {(assets ?? [])
                    .filter((a) => a.type === "video")
                    .map((a) => (
                      <Paper key={a.id} p="sm" radius="sm" className="bg-studio-card border border-studio-border">
                        <Group justify="space-between">
                          <Text size="xs" c="white" className="font-mono">{a.id} — {a.status}</Text>
                          <Text size="xs" c="dimmed">{a.provider}/{a.model}</Text>
                        </Group>
                      </Paper>
                    ))}
                </Stack>
              ) : (
                <Stack align="center" gap="sm" py="xl">
                  <ThemeIcon variant="light" color="terracotta" size={40} radius="md">
                    <Video size={20} />
                  </ThemeIcon>
                  <Text size="xs" c="dimmed">
                    Generated video shots with 9:16 vertical previews
                  </Text>
                </Stack>
              )}
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
              className="bg-studio-panel border border-studio-border min-h-[200px]"
            >
              {(assets ?? []).length > 0 ? (
                <Stack gap="xs">
                  {(assets ?? []).map((a) => (
                    <Paper key={a.id} p="sm" radius="sm" className="bg-studio-card border border-studio-border">
                      <Group justify="space-between">
                        <Text size="xs" c="white" className="font-mono">{a.type} — {a.id}</Text>
                        <Badge color="gray" variant="outline" size="xs">{a.status}</Badge>
                      </Group>
                    </Paper>
                  ))}
                </Stack>
              ) : (
                <Text size="xs" c="dimmed" ta="center" py="xl">
                  All generated media assets for this episode — images, video,
                  audio, renders.
                </Text>
              )}
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
              <Button
                variant="outline"
                color="terracotta"
                size="xs"
                loading={generateAsset.isPending}
                onClick={() =>
                  generateAsset.mutate({
                    capability: "voice",
                    type: "voice",
                    prompt: `Dialogue voiceover for ${episode?.title ?? "episode"}`,
                  })
                }
              >
                Generate Audio
              </Button>
            </Group>
            <Stack gap="xs">
              {(["voice", "music", "sfx"] as const).map((t) => {
                const trackAssets = (assets ?? []).filter((a) => a.type === t);
                const label = t === "voice" ? "Dialogue" : t === "music" ? "Music / Score" : "Sound Effects";
                return (
                  <Paper
                    key={t}
                    p="sm"
                    radius="sm"
                    className="bg-studio-panel border border-studio-border"
                  >
                    <Group justify="space-between">
                      <Text size="xs" fw={600} c="white">
                        {label}
                      </Text>
                      <Badge color={trackAssets.length > 0 ? "emerald" : "gray"} variant="outline" size="xs">
                        {trackAssets.length > 0 ? `${trackAssets.length} generated` : "Not generated"}
                      </Badge>
                    </Group>
                  </Paper>
                );
              })}
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
              <Button
                variant="filled"
                color="terracotta"
                size="xs"
                loading={renderEpisode.isPending}
                onClick={() => renderEpisode.mutate()}
              >
                Render Episode
              </Button>
            </Group>
            <Stack gap="sm">
              {[
                { label: "Video Track", items: timeline?.video_tracks ?? [] },
                { label: "Audio Track", items: timeline?.audio_tracks ?? [] },
              ].map((track) => (
                <div key={track.label} className="space-y-1">
                  <Text size="xs" c="dimmed" fw={600}>
                    {track.label}
                    {timeline ? ` — v${timeline.version} (${timeline.status})` : ""}
                  </Text>
                  <div className="h-10 bg-studio-panel border border-studio-border rounded p-2 flex gap-2 overflow-x-auto">
                    {track.items.length > 0 ? (
                      track.items.map((t) => (
                        <div key={t.id} className="h-full bg-studio-accent/10 border border-studio-accent/20 rounded px-3 flex items-center text-xs text-studio-muted font-mono min-w-[160px]">
                          {t.shot_id} · {t.start_time}s +{t.duration}s
                        </div>
                      ))
                    ) : (
                      <div className="h-full bg-studio-accent/10 border border-studio-accent/20 rounded px-3 flex items-center text-xs text-studio-muted font-mono min-w-[160px]">
                        No clips
                      </div>
                    )}
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
              className="bg-studio-panel border border-studio-border min-h-[200px]"
            >
              {(continuityIssues ?? []).length > 0 ? (
                <Stack gap="xs">
                  {(continuityIssues ?? []).map((issue) => (
                    <Paper key={issue.id} p="sm" radius="sm" className="bg-studio-card border border-studio-border">
                      <Group justify="space-between">
                        <Text size="xs" c="white">{issue.entity}: expected {issue.expected_state}, found {issue.actual_state}</Text>
                        <Badge
                          color={issue.status === "resolved" || issue.status === "wontfix" ? "emerald" : "amber"}
                          variant="light"
                          size="xs"
                        >
                          {issue.status}
                        </Badge>
                      </Group>
                    </Paper>
                  ))}
                </Stack>
              ) : (
                <Stack align="center" gap="sm" py="xl">
                  <ThemeIcon variant="light" color="emerald" size={40} radius="md">
                    <ShieldCheck size={20} />
                  </ThemeIcon>
                  <Text size="xs" c="dimmed">
                    Continuity checks run automatically. Issues appear here when
                    detected.
                  </Text>
                </Stack>
              )}
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
