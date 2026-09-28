"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Activity, CheckCircle2, XCircle, RotateCcw, Shield, Cpu } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, Button, ThemeIcon, Alert, SimpleGrid, Select } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { LeadDirectorDecision, ApprovalRequest, Season, Episode, ProductionRun } from "@/lib/api/types";

export default function ProductionControlTowerPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const queryClient = useQueryClient();

  const { data: decisions } = useQuery({
    queryKey: ["decisions", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: LeadDirectorDecision[] }>(`/v1/projects/${projectId}/agents/decisions`);
      return res.items || [];
    },
  });

  const { data: approvals } = useQuery({
    queryKey: ["approvals", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: ApprovalRequest[] }>(`/v1/projects/${projectId}/production/approvals`);
      return res.items || [];
    },
  });

  const { data: runs } = useQuery({
    queryKey: ["production-runs", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: ProductionRun[] }>(`/v1/projects/${projectId}/production/runs`);
      return res.items || [];
    },
  });

  // Episodes live under seasons — flatten all seasons' episodes into one
  // selector so director steps and approvals target a real episode.
  const { data: seasons } = useQuery({
    queryKey: ["seasons", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: Season[] }>(`/v1/projects/${projectId}/seasons`);
      return res.items || [];
    },
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
    enabled: Boolean(seasons && seasons.length > 0),
  });

  const [episodeId, setEpisodeId] = useState("");
  useEffect(() => {
    if (!episodeId && episodes && episodes.length > 0) {
      setEpisodeId(episodes[0].id);
    }
  }, [episodes, episodeId]);

  const triggerDirectorMutation = useMutation({
    mutationFn: async () => {
      return api.post<{ result: Record<string, unknown>; error?: string }>(
        `/v1/projects/${projectId}/agents/director/step`,
        { episode_id: episodeId }
      );
    },
    onSuccess: (d) => {
      queryClient.invalidateQueries({ queryKey: ["decisions", projectId] });
      notifications.show({
        title: "Lead Director Loop Executed",
        message: d.error ? `Step error: ${d.error}` : "Director step completed",
        color: "terracotta",
      });
    },
  });

  const approvalMutation = useMutation({
    mutationFn: async (decision: "APPROVE" | "REJECT" | "REQUEST_REVISION") => {
      const req = await api.post<ApprovalRequest>(`/v1/projects/${projectId}/production/approvals`, {
        episode_id: episodeId,
        stage: "Episode",
        target_id: episodeId,
      });
      return api.post<ApprovalRequest>(
        `/v1/projects/${projectId}/production/approvals/${req.id}/decide`,
        { decision, notes: "Decision submitted via Production Console" }
      );
    },
    onSuccess: (a) => {
      queryClient.invalidateQueries({ queryKey: ["approvals", projectId] });
      notifications.show({
        title: "Human Gate Submitted",
        message: `Status set to ${a.decision}`,
        color: a.decision === "APPROVE" ? "emerald" : "amber",
      });
    },
  });

  const runAction = useMutation({
    mutationFn: async ({ runId, action }: { runId: string; action: "pause" | "resume" | "stop" }) =>
      api.post<ProductionRun>(`/v1/projects/${projectId}/production/runs/${runId}/${action}`, {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["production-runs", projectId] });
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="terracotta" variant="light" size={44} radius="md">
              <Activity size={22} />
            </ThemeIcon>
            <div>
              <Title order={3} c="white">
                Production Console
              </Title>
              <Text size="xs" c="dimmed">
                Lead Director console, active production workflows, and human approval gates
              </Text>
            </div>
          </Group>

          <Group gap="sm">
            <Select
              size="xs"
              placeholder="Select episode"
              value={episodeId || null}
              onChange={(v) => setEpisodeId(v || "")}
              data={(episodes || []).map((e) => ({
                value: e.id,
                label: `E${String(e.number).padStart(2, "0")} — ${e.title}`,
              }))}
              styles={{ input: { backgroundColor: "#1a1a1a", color: "#fff", borderColor: "#333" } }}
            />
            <Button
              onClick={() => triggerDirectorMutation.mutate()}
              loading={triggerDirectorMutation.isPending}
              disabled={!episodeId}
              leftSection={<Cpu size={16} />}
              variant="filled"
              color="terracotta"
              size="sm"
              radius="sm"
            >
              Execute Lead Director Step
            </Button>
          </Group>
        </Group>
      </Paper>

      {/* Human Approval Gate Box */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-amber-500/30">
        <Stack gap="md">
          <Group justify="space-between">
            <Group gap="xs">
              <span className="w-2.5 h-2.5 rounded-full bg-amber-400"></span>
              <Text fw={700} size="sm" c="amber.4" tt="uppercase" className="tracking-wider font-mono">
                Human Approval Gate
              </Text>
            </Group>
            <Badge color="amber" variant="light" size="sm" className="font-mono">
              Production Workflow Paused
            </Badge>
          </Group>

          <Text size="xs" c="white">
            Episode 01 Storyboard & Shot Generations require human approval before advancing to Assembly.
          </Text>

          <Group gap="sm" pt="xs">
            <Button
              onClick={() => approvalMutation.mutate("APPROVE")}
              disabled={!episodeId}
              loading={approvalMutation.isPending}
              leftSection={<CheckCircle2 size={16} />}
              color="emerald"
              size="xs"
              radius="sm"
            >
              Approve Production
            </Button>
            <Button
              onClick={() => approvalMutation.mutate("REQUEST_REVISION")}
              loading={approvalMutation.isPending}
              leftSection={<RotateCcw size={16} />}
              variant="light"
              color="amber"
              size="xs"
              radius="sm"
            >
              Request Revision
            </Button>
            <Button
              onClick={() => approvalMutation.mutate("REJECT")}
              loading={approvalMutation.isPending}
              leftSection={<XCircle size={16} />}
              variant="light"
              color="red"
              size="xs"
              radius="sm"
            >
              Reject & Regenerate
            </Button>
          </Group>
        </Stack>
      </Paper>

      {/* Active Production Runs */}
      <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Stack gap="md">
          <Group gap="xs">
            <Activity size={16} className="text-studio-accent" />
            <Text fw={700} size="sm" c="white">
              Production Runs
            </Text>
          </Group>
          <Stack gap="xs">
            {runs && runs.length > 0 ? (
              runs.map((run) => (
                <Paper key={run.id} p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
                  <Group justify="space-between">
                    <Group gap="xs">
                      <Text fw={700} size="xs" c="white" className="font-mono">
                        {run.id}
                      </Text>
                      <Badge size="xs" variant="light" color={run.status === "running" ? "emerald" : run.status === "paused" ? "amber" : "gray"}>
                        {run.status}
                      </Badge>
                    </Group>
                    <Group gap="xs">
                      {run.status === "running" && (
                        <Button size="xs" variant="light" color="amber"
                          onClick={() => runAction.mutate({ runId: run.id, action: "pause" })}>
                          Pause
                        </Button>
                      )}
                      {run.status === "paused" && (
                        <Button size="xs" variant="light" color="emerald"
                          onClick={() => runAction.mutate({ runId: run.id, action: "resume" })}>
                          Resume
                        </Button>
                      )}
                      {(run.status === "running" || run.status === "paused") && (
                        <Button size="xs" variant="light" color="red"
                          onClick={() => runAction.mutate({ runId: run.id, action: "stop" })}>
                          Stop
                        </Button>
                      )}
                    </Group>
                  </Group>
                  <Text size="xs" c="dimmed" mt={4}>
                    Episode {run.episode_id} — stage {run.stage}
                  </Text>
                </Paper>
              ))
            ) : (
              <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">
                No production runs yet.
              </Text>
            )}
          </Stack>
        </Stack>
      </Paper>

      {/* Logs & Decision Grid */}
      <SimpleGrid cols={{ base: 1, md: 2 }} spacing="lg">
        {/* Lead Director Decision Records */}
        <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
          <Stack gap="md">
            <Group gap="xs">
              <Cpu size={16} className="text-studio-accent" />
              <Text fw={700} size="sm" c="white">
                Auditable Lead Director Decisions
              </Text>
            </Group>

            <Stack gap="xs">
              {decisions && decisions.length > 0 ? (
                decisions.map((d) => (
                  <Paper key={d.id} p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
                    <Group justify="space-between">
                      <Text fw={700} size="xs" c="terracotta.4" className="font-mono">
                        {d.decision}
                      </Text>
                      <Text size="xs" c="dimmed" className="font-mono">
                        {d.decision_maker}
                      </Text>
                    </Group>
                    <Text size="xs" c="white" mt={4}>
                      {d.reason}
                    </Text>
                  </Paper>
                ))
              ) : (
                <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">
                  No decisions recorded yet. Click "Execute Lead Director Step" above.
                </Text>
              )}
            </Stack>
          </Stack>
        </Paper>

        {/* Approval History */}
        <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
          <Stack gap="md">
            <Group gap="xs">
              <Shield size={16} className="text-emerald-400" />
              <Text fw={700} size="sm" c="white">
                Approval History
              </Text>
            </Group>

            <Stack gap="xs">
              {approvals && approvals.length > 0 ? (
                approvals.map((a) => (
                  <Paper key={a.id} p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
                    <Group justify="space-between">
                      <Badge color={a.decision === "APPROVE" ? "emerald" : "amber"} variant="light" size="xs" className="font-mono">
                        {a.decision}
                      </Badge>
                      <Text size="xs" c="dimmed" className="font-mono">
                        {a.decided_by}
                      </Text>
                    </Group>
                    <Text size="xs" c="white" mt={4}>
                      {a.notes}
                    </Text>
                  </Paper>
                ))
              ) : (
                <Text size="xs" c="dimmed" fs="italic" ta="center" py="md">
                  No human approvals recorded yet.
                </Text>
              )}
            </Stack>
          </Stack>
        </Paper>
      </SimpleGrid>
    </Stack>
  );
}
