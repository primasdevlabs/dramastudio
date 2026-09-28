"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Activity, CheckCircle2, XCircle, RotateCcw, Shield, Cpu } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, Button, ThemeIcon, Alert, SimpleGrid } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { LeadDirectorDecision, ApprovalRequest } from "@/lib/api/types";

export default function ProductionControlTowerPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const { data: decisions } = useQuery({
    queryKey: ["decisions", projectId],
    queryFn: async () => {
      const res = await api.get<{ decisions: LeadDirectorDecision[] }>(`/v1/agents/decisions?project_id=${projectId}`);
      return res.decisions || [];
    },
  });

  const { data: approvals } = useQuery({
    queryKey: ["approvals", projectId],
    queryFn: async () => {
      const res = await api.get<{ approvals: ApprovalRequest[] }>(`/v1/production/approval?project_id=${projectId}`);
      return res.approvals || [];
    },
  });

  const triggerDirectorMutation = useMutation({
    mutationFn: async () => {
      return api.post<LeadDirectorDecision>("/v1/agents/lead-director/step", {
        project_id: projectId,
        episode_id: "ep_001",
      });
    },
    onSuccess: (d) => {
      queryClient.invalidateQueries({ queryKey: ["decisions", projectId] });
      notifications.show({
        title: "Lead Director Loop Executed",
        message: `Decision: ${d.decision || "Step Completed"}`,
        color: "terracotta",
      });
    },
  });

  const approvalMutation = useMutation({
    mutationFn: async (decision: "APPROVE" | "REJECT" | "REQUEST_REVISION") => {
      return api.post<ApprovalRequest>("/v1/production/approval", {
        project_id: projectId,
        episode_id: "ep_001",
        stage: "Episode",
        target_id: "ep_001",
        decision,
        notes: "Decision submitted via Production Control Tower",
        decided_by: "Lead Producer",
      });
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

          <Button
            onClick={() => triggerDirectorMutation.mutate()}
            loading={triggerDirectorMutation.isPending}
            leftSection={<Cpu size={16} />}
            variant="filled"
            color="terracotta"
            size="sm"
            radius="sm"
          >
            Execute Lead Director Step
          </Button>
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
