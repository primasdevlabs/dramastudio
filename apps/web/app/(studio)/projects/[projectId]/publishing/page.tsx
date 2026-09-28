"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import {
  Share2,
  Send,
  CheckCircle2,
  Calendar as CalendarIcon,
  Film,
  Video,
  Play,
  AlertTriangle,
  RefreshCw,
  Edit3,
  Layers,
  Clock,
  ExternalLink,
  Check,
  X,
  FileVideo,
  Info,
  Sliders,
  CheckSquare,
  Square,
  HelpCircle,
  ArrowRight,
} from "lucide-react";
import {
  Paper,
  Title,
  Text,
  Badge,
  Button,
  Tabs,
  Group,
  Stack,
  Divider,
  Modal,
  Switch,
  Select,
  TextInput,
  Textarea,
  Alert,
  Tooltip,
} from "@mantine/core";
import { useSettingsStore } from "@/stores/settings-store";
import type { PublishingPlan, PublishingPlanChannelItem } from "@/lib/api/settings-types";
import { api } from "@/lib/api/client";
import type { Channel, Publication, Season, Episode } from "@/lib/api/types";

export default function PublishingPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const queryClient = useQueryClient();
  const { getEffectiveSettings } = useSettingsStore();
  const effectiveSettings = getEffectiveSettings(projectId);

  const [activeTab, setActiveTab] = useState("overview");

  // Publishing Plan Modal state
  const [isPlanModalOpen, setIsPlanModalOpen] = useState(false);
  const [editingPlan, setEditingPlan] = useState<PublishingPlan | null>(null);

  // Failure Detail Modal state
  const [failureDetail, setFailureDetail] = useState<{ channel: string; reason: string } | null>(null);

  const { data: channels } = useQuery({
    queryKey: ["channels", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: Channel[] }>(
        `/v1/projects/${projectId}/publishing/channels`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });

  const { data: publications } = useQuery({
    queryKey: ["publications", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: Publication[] }>(
        `/v1/projects/${projectId}/publishing/publications`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });

  const { data: seasons } = useQuery({
    queryKey: ["seasons", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: Season[] }>(
        `/v1/projects/${projectId}/seasons`
      );
      return res.items ?? [];
    },
    enabled: Boolean(projectId),
  });

  const { data: episodes } = useQuery({
    queryKey: ["project-episodes", projectId, (seasons || []).map((s) => s.id).join(",")],
    queryFn: async () => {
      const all: Episode[] = [];
      for (const s of seasons || []) {
        const res = await api.get<{ items: Episode[] }>(
          `/v1/projects/${projectId}/seasons/${s.id}/episodes`
        );
        all.push(...(res.items || []));
      }
      return all;
    },
    enabled: Boolean(projectId && seasons && seasons.length > 0),
  });

  const episodeTitle = (id: string) => {
    const ep = (episodes || []).find((e) => e.id === id);
    return ep ? `E${String(ep.number).padStart(2, "0")} — ${ep.title}` : id;
  };

  // Group real publications into per-episode plans; connected channels not
  // yet scheduled for that episode appear disabled.
  const plans: PublishingPlan[] = [];
  const byEpisode = new Map<string, Publication[]>();
  for (const pub of publications ?? []) {
    byEpisode.set(pub.episode_id, [...(byEpisode.get(pub.episode_id) ?? []), pub]);
  }
  for (const [episodeId, pubs] of byEpisode) {
    const scheduled = pubs.find((p) => p.scheduled_at)?.scheduled_at;
    plans.push({
      id: `plan-${episodeId}`,
      episodeId,
      episodeTitle: episodeTitle(episodeId),
      scheduledTime: scheduled ? new Date(scheduled).toLocaleString() : "Immediately",
      channels: (channels ?? []).map((ch) => {
        const pub = pubs.find((p) => p.channel_id === ch.id);
        return {
          channelId: ch.id,
          channelName: ch.name || ch.platform,
          enabled: Boolean(pub),
          format: pub?.metadata?.title ?? "Video",
          publicationTime: pub?.scheduled_at
            ? new Date(pub.scheduled_at).toLocaleString()
            : "Immediately",
          status: (pub?.status ?? "draft") as PublishingPlanChannelItem["status"],
        };
      }),
    });
  }

  const scheduleMutation = useMutation({
    mutationFn: async (plan: PublishingPlan) => {
      const existing = new Set(
        (publications ?? [])
          .filter((p) => p.episode_id === plan.episodeId)
          .map((p) => p.channel_id)
      );
      for (const ch of plan.channels) {
        if (ch.enabled && !existing.has(ch.channelId)) {
          await api.post(`/v1/projects/${projectId}/publishing/publications`, {
            episode_id: plan.episodeId,
            channel_id: ch.channelId,
            metadata: { title: plan.episodeTitle, caption: "", tags: [] },
          });
        }
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["publications", projectId] });
    },
  });

  const handleOpenEditPlan = (plan: PublishingPlan) => {
    setEditingPlan(JSON.parse(JSON.stringify(plan)));
    setIsPlanModalOpen(true);
  };

  const handleSavePlan = () => {
    if (!editingPlan) return;
    scheduleMutation.mutate(editingPlan);
    setIsPlanModalOpen(false);
  };

  const statusBadges: Record<string, { color: string; label: string }> = {
    scheduled: { color: "blue", label: "Scheduled" },
    publishing: { color: "cyan", label: "Publishing" },
    published: { color: "emerald", label: "Published" },
    failed: { color: "red", label: "Failed" },
    draft: { color: "gray", label: "Draft" },
  };

  return (
    <div className="max-w-7xl mx-auto space-y-6 pb-12">
      {/* Workspace Header */}
      <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded bg-studio-accent/10 border border-studio-accent/20 flex items-center justify-center text-studio-accent">
              <Share2 className="w-5 h-5" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <Title order={2} className="text-xl font-bold text-white tracking-tight">
                  Publishing Workspace
                </Title>
                <Badge color="terracotta" variant="light" size="sm">
                  PUBLISHING
                </Badge>
              </div>
              <Text size="xs" c="dimmed" mt={2}>
                Decide what gets published, where, and when across connected distribution channels.
              </Text>
            </div>
          </div>

          <Button color="terracotta" size="xs" leftSection={<Send size={14} />}>
            Create Distribution Package
          </Button>
        </div>
      </Paper>

      {/* Workspace Navigation Tabs */}
      <Tabs
        value={activeTab}
        onChange={(val) => val && setActiveTab(val)}
        variant="outline"
        classNames={{ list: "border-studio-border" }}
      >
        <Tabs.List className="mb-6 flex-wrap">
          <Tabs.Tab value="overview" leftSection={<Layers size={14} />}>
            Overview & Deliverables
          </Tabs.Tab>
          <Tabs.Tab value="calendar" leftSection={<CalendarIcon size={14} />}>
            Calendar
          </Tabs.Tab>
          <Tabs.Tab value="episodes" leftSection={<Film size={14} />}>
            Episodes & Publishing Plans
          </Tabs.Tab>
          <Tabs.Tab value="clips" leftSection={<Video size={14} />}>
            Clips & Shorts
          </Tabs.Tab>
          <Tabs.Tab value="channels" leftSection={<Share2 size={14} />}>
            Connected Channels
          </Tabs.Tab>
          <Tabs.Tab value="history" leftSection={<Clock size={14} />}>
            Publication History
          </Tabs.Tab>
        </Tabs.List>

        {/* Tab 1: Overview & Deliverables Architecture */}
        <Tabs.Panel value="overview">
          <div className="space-y-6">
            {/* Master Asset Pipeline Card */}
            <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
              <div>
                <Text fw={700} size="sm" c="white" className="flex items-center gap-2 uppercase text-xs tracking-wider text-studio-accent">
                  <FileVideo size={16} />
                  Distribution Deliverables
                </Text>
                <Text size="xs" c="dimmed">
                  Master episodes adapted for each connected channel.
                </Text>
              </div>

              <Divider color="dark.5" />

              <div className="bg-studio-panel p-6 rounded border border-studio-border">
                <div className="grid grid-cols-1 md:grid-cols-5 gap-4 items-center text-center">
                  {/* Master Render */}
                  <div className="p-4 bg-studio-card border border-studio-accent/40 rounded flex flex-col items-center justify-center space-y-2">
                    <Film className="w-8 h-8 text-studio-accent" />
                    <Text fw={700} size="xs" c="white">Episode Master</Text>
                    <Text size="xs" c="dimmed">Master Video</Text>
                  </div>

                  <div className="hidden md:flex justify-center text-studio-muted"><ArrowRight size={20} /></div>

                  {/* Deliverables Breakdown */}
                  <div className="md:col-span-3 grid grid-cols-2 sm:grid-cols-4 gap-3 text-left">
                    <div className="p-3 bg-studio-card border border-studio-border rounded space-y-1">
                      <div className="flex items-center gap-1.5 text-red-400">
                        <Play size={14} />
                        <Text fw={700} size="xs">YouTube</Text>
                      </div>
                      <Text size="xs" c="dimmed">Full Episode (16:9) / Shorts (9:16)</Text>
                    </div>

                    <div className="p-3 bg-studio-card border border-studio-border rounded space-y-1">
                      <div className="flex items-center gap-1.5 text-gray-300">
                        <Video size={14} />
                        <Text fw={700} size="xs">TikTok</Text>
                      </div>
                      <Text size="xs" c="dimmed">Vertical Episode (9:16)</Text>
                    </div>

                    <div className="p-3 bg-studio-card border border-studio-border rounded space-y-1">
                      <div className="flex items-center gap-1.5 text-purple-400">
                        <Video size={14} />
                        <Text fw={700} size="xs">Instagram</Text>
                      </div>
                      <Text size="xs" c="dimmed">Reel (9:16)</Text>
                    </div>

                    <div className="p-3 bg-studio-card border border-studio-border rounded space-y-1">
                      <div className="flex items-center gap-1.5 text-blue-400">
                        <Video size={14} />
                        <Text fw={700} size="xs">Facebook</Text>
                      </div>
                      <Text size="xs" c="dimmed">Video / Reel</Text>
                    </div>
                  </div>
                </div>
              </div>
            </Paper>

            {/* Current Active Publishing Plans */}
            <div className="space-y-4">
              <Text fw={700} size="sm" c="white" className="uppercase text-xs tracking-wider text-studio-accent">
                Active Publishing Plans
              </Text>

              {plans.length === 0 && (
                <Paper p="lg" radius="md" className="bg-studio-card border border-studio-border">
                  <Text size="xs" c="dimmed" fs="italic" ta="center">
                    No publications scheduled yet. Edit a publishing plan to schedule one.
                  </Text>
                </Paper>
              )}

              {plans.map((plan) => (
                <Paper key={plan.id} p="lg" radius="md" className="bg-studio-card border border-studio-border space-y-4">
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <div>
                      <Text fw={700} size="sm" c="white">{plan.episodeTitle}</Text>
                      <Text size="xs" c="dimmed" mt={1}>
                        Target Schedule: <span className="text-white">{plan.scheduledTime}</span>
                      </Text>
                    </div>

                    <Button
                      variant="outline"
                      color="terracotta"
                      size="xs"
                      leftSection={<Edit3 size={14} />}
                      onClick={() => handleOpenEditPlan(plan)}
                    >
                      Edit Publishing Plan
                    </Button>
                  </div>

                  <Divider color="dark.5" />

                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
                    {plan.channels.map((ch) => (
                      <Paper
                        key={ch.channelId}
                        p="sm"
                        radius="sm"
                        className="bg-studio-panel border border-studio-border space-y-2"
                      >
                        <div className="flex items-center justify-between">
                          <Text fw={700} size="xs" c="white">{ch.channelName}</Text>
                          <Badge
                            color={statusBadges[ch.status]?.color}
                            variant="light"
                            size="xs"
                          >
                            ● {statusBadges[ch.status]?.label}
                          </Badge>
                        </div>

                        <div className="text-xs space-y-1">
                          <Text size="xs" c="dimmed">Format: <span className="text-white">{ch.format}</span></Text>
                          <Text size="xs" c="dimmed">Time: <span className="text-white">{ch.publicationTime}</span></Text>
                        </div>

                        {ch.status === "failed" && (
                          <div className="pt-2 border-t border-studio-border flex justify-between items-center text-xs">
                            <span className="text-red-400 font-semibold flex items-center gap-1">
                              <AlertTriangle size={12} /> Failed
                            </span>
                            <Button
                              variant="subtle"
                              color="red"
                              size="compact-xs"
                              onClick={() =>
                                setFailureDetail({
                                  channel: ch.channelName,
                                  reason: ch.failureReason || "Publication failed.",
                                })
                              }
                            >
                              View Details
                            </Button>
                          </div>
                        )}
                      </Paper>
                    ))}
                  </div>
                </Paper>
              ))}
            </div>
          </div>
        </Tabs.Panel>

        {/* Tab 2: Calendar */}
        <Tabs.Panel value="calendar">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border text-center space-y-3">
            <CalendarIcon size={32} className="mx-auto text-studio-accent" />
            <Text fw={700} size="sm" c="white">Distribution Calendar</Text>
            <Text size="xs" c="dimmed" className="max-w-md mx-auto">
              Visual release schedule showing publication windows across connected channels.
            </Text>
          </Paper>
        </Tabs.Panel>

        {/* Tab 3: Episodes & Plans */}
        <Tabs.Panel value="episodes">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-4">
            <Text fw={700} size="sm" c="white">Episode Publishing Plans</Text>
            <Text size="xs" c="dimmed">Manage publishing plans and captions for each finished episode.</Text>
          </Paper>
        </Tabs.Panel>

        {/* Tab 4: Clips & Shorts */}
        <Tabs.Panel value="clips">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-4">
            <Text fw={700} size="sm" c="white">Clips & Shorts</Text>
            <Text size="xs" c="dimmed">Distribute short-form clips extracted from master episodes.</Text>
          </Paper>
        </Tabs.Panel>

        {/* Tab 5: Connected Channels Capability View */}
        <Tabs.Panel value="channels">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
            <div>
              <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                <Share2 size={16} className="text-studio-accent" />
                Channel Capabilities
              </Text>
              <Text size="xs" c="dimmed">
                Supported features and formats for each connected channel.
              </Text>
            </div>

            <Divider color="dark.5" />

            {(channels ?? []).length > 0 && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {(channels ?? []).map((ch) => (
                  <Paper key={ch.id} p="md" radius="sm" className="bg-studio-panel border border-studio-border space-y-3">
                    <Group justify="space-between">
                      <Text fw={700} size="sm" c="white">{ch.name || ch.platform}</Text>
                      <Badge color={ch.enabled ? "emerald" : "gray"} variant="light" size="xs">
                        {ch.enabled ? "CONNECTED" : "DISABLED"}
                      </Badge>
                    </Group>
                    <Text size="xs" c="dimmed">Platform: {ch.platform}{ch.account_ref ? ` — ${ch.account_ref}` : ""}</Text>
                  </Paper>
                ))}
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Facebook */}
              <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border space-y-3">
                <Group justify="space-between">
                  <Text fw={700} size="sm" c="white">Facebook</Text>
                  <Badge color={effectiveSettings.facebook.connected ? "emerald" : "gray"} variant="light" size="xs">
                    {effectiveSettings.facebook.connected ? "CONNECTED" : "DISCONNECTED"}
                  </Badge>
                </Group>
                <Text size="xs" c="dimmed">Capabilities: Video, Image, Scheduling</Text>
              </Paper>

              {/* Instagram */}
              <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border space-y-3">
                <Group justify="space-between">
                  <Text fw={700} size="sm" c="white">Instagram</Text>
                  <Badge color={effectiveSettings.instagram.connected ? "emerald" : "gray"} variant="light" size="xs">
                    {effectiveSettings.instagram.connected ? "CONNECTED" : "DISCONNECTED"}
                  </Badge>
                </Group>
                <Text size="xs" c="dimmed">Capabilities: Reels (9:16), Image, Carousel, Scheduling</Text>
              </Paper>

              {/* TikTok */}
              <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border space-y-3">
                <Group justify="space-between">
                  <Text fw={700} size="sm" c="white">TikTok</Text>
                  <Badge color={effectiveSettings.tiktok.connected ? "emerald" : "gray"} variant="light" size="xs">
                    {effectiveSettings.tiktok.connected ? "CONNECTED" : "DISCONNECTED"}
                  </Badge>
                </Group>
                <Text size="xs" c="dimmed">Capabilities: Vertical Video (9:16), Scheduling</Text>
              </Paper>

              {/* YouTube */}
              <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border space-y-3">
                <Group justify="space-between">
                  <Text fw={700} size="sm" c="white">YouTube</Text>
                  <Badge color={effectiveSettings.youtube.connected ? "emerald" : "gray"} variant="light" size="xs">
                    {effectiveSettings.youtube.connected ? "CONNECTED" : "DISCONNECTED"}
                  </Badge>
                </Group>
                <Text size="xs" c="dimmed">Capabilities: Full Video (16:9), Shorts (9:16), Thumbnails, Scheduling</Text>
              </Paper>
            </div>
          </Paper>
        </Tabs.Panel>

        {/* Tab 6: History */}
        <Tabs.Panel value="history">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-4">
            <Text fw={700} size="sm" c="white">Publication History</Text>
            <Text size="xs" c="dimmed">History of all published episodes and clips across your channels.</Text>
            {(publications ?? []).length === 0 ? (
              <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">No publications yet.</Text>
            ) : (
              <Stack gap="xs">
                {(publications ?? []).map((pub) => (
                  <Paper key={pub.id} p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
                    <Group justify="space-between">
                      <div>
                        <Text fw={700} size="xs" c="white">{pub.metadata?.title || episodeTitle(pub.episode_id)}</Text>
                        <Text size="xs" c="dimmed" className="font-mono">
                          {pub.channel_id} · {pub.scheduled_at ? new Date(pub.scheduled_at).toLocaleString() : "immediate"}
                        </Text>
                      </div>
                      <Badge color={statusBadges[pub.status]?.color ?? "gray"} variant="light" size="xs">
                        {statusBadges[pub.status]?.label ?? pub.status}
                      </Badge>
                    </Group>
                  </Paper>
                ))}
              </Stack>
            )}
          </Paper>
        </Tabs.Panel>
      </Tabs>

      {/* PUBLISHING PLAN EDIT MODAL */}
      <Modal
        opened={isPlanModalOpen}
        onClose={() => setIsPlanModalOpen(false)}
        title={
          <Text fw={700} size="sm" c="white">
            PUBLISHING PLAN — {editingPlan?.episodeTitle}
          </Text>
        }
        centered
        size="lg"
        classNames={{
          content: "bg-studio-card border border-studio-border text-white",
          header: "bg-studio-card border-b border-studio-border text-white",
        }}
      >
        {editingPlan && (
          <div className="space-y-6 pt-2">
            <Text size="xs" c="dimmed">
              Select channels, formats, and scheduled release times for this episode.
            </Text>

            <div className="space-y-4">
              {editingPlan.channels.map((ch, idx) => (
                <Paper
                  key={ch.channelId}
                  p="md"
                  radius="sm"
                  className="bg-studio-panel border border-studio-border space-y-3"
                >
                  <Group justify="space-between">
                    <Group gap="sm">
                      <Switch
                        color="terracotta"
                        checked={ch.enabled}
                        onChange={(e) => {
                          const updatedChannels = [...editingPlan.channels];
                          updatedChannels[idx].enabled = e.currentTarget.checked;
                          setEditingPlan({ ...editingPlan, channels: updatedChannels });
                        }}
                      />
                      <Text fw={700} size="xs" c="white">{ch.channelName}</Text>
                    </Group>
                    <Badge color={ch.enabled ? "terracotta" : "gray"} variant="light" size="xs">
                      {ch.enabled ? "ENABLED" : "EXCLUDED"}
                    </Badge>
                  </Group>

                  {ch.enabled && (
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
                      <Select
                        label="Format"
                        size="xs"
                        value={ch.format}
                        onChange={(val) => {
                          if (!val) return;
                          const updatedChannels = [...editingPlan.channels];
                          updatedChannels[idx].format = val;
                          setEditingPlan({ ...editingPlan, channels: updatedChannels });
                        }}
                        data={["Video (16:9)", "Reel (9:16)", "Shorts (9:16)", "Full Episode (16:9)"]}
                      />

                      <TextInput
                        label="Publication Time"
                        size="xs"
                        value={ch.publicationTime}
                        onChange={(e) => {
                          const updatedChannels = [...editingPlan.channels];
                          updatedChannels[idx].publicationTime = e.currentTarget.value;
                          setEditingPlan({ ...editingPlan, channels: updatedChannels });
                        }}
                      />
                    </div>
                  )}
                </Paper>
              ))}
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t border-studio-border">
              <Button variant="outline" color="gray" size="xs" onClick={() => setIsPlanModalOpen(false)}>
                Cancel
              </Button>
              <Button color="terracotta" size="xs" onClick={handleSavePlan}>
                Save Plan
              </Button>
            </div>
          </div>
        )}
      </Modal>

      {/* DIAGNOSTIC FAILURE DETAIL MODAL */}
      <Modal
        opened={!!failureDetail}
        onClose={() => setFailureDetail(null)}
        title={
          <div className="flex items-center gap-2 text-red-400">
            <AlertTriangle size={18} />
            <Text fw={700} size="sm" c="white">
              Publication Issue — {failureDetail?.channel}
            </Text>
          </div>
        }
        centered
        size="md"
        classNames={{
          content: "bg-studio-card border border-studio-border text-white",
          header: "bg-studio-card border-b border-studio-border text-white",
        }}
      >
        {failureDetail && (
          <div className="space-y-4 pt-2">
            <Alert color="red" variant="outline" icon={<AlertTriangle size={16} />}>
              <Text size="xs" fw={700} c="red">Publication Failed</Text>
              <Text size="xs" c="dimmed" mt={1}>
                {failureDetail.reason}
              </Text>
            </Alert>

            <div className="bg-studio-panel p-4 rounded border border-studio-border space-y-1 text-xs">
              <Text fw={600} c="white">Next Step:</Text>
              <p className="text-studio-muted">
                Ensure the master asset format matches the channel requirement and try publishing again.
              </p>
            </div>

            <div className="flex justify-end gap-3 pt-3 border-t border-studio-border">
              <Button variant="outline" color="gray" size="xs" onClick={() => setFailureDetail(null)}>
                Close
              </Button>
              <Button color="terracotta" size="xs" leftSection={<RefreshCw size={14} />}>
                Retry Publication
              </Button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
