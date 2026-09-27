"use client";

import { Paper, Group, Stack, Text, Badge, ThemeIcon, Progress } from "@mantine/core";
import { Cpu, Activity, Eye, CheckCircle2, ListChecks } from "lucide-react";

export function LeadDirectorConsole({
  phase = "Plan & Assess",
  stepCount = 4,
  status = "Active",
}: {
  phase?: string;
  stepCount?: number;
  status?: string;
}) {
  return (
    <Paper p="lg" radius="xl" withBorder className="bg-studio-card border-studio-border shadow-xl">
      <Stack gap="md">
        <Group justify="space-between">
          <Group gap="md">
            <ThemeIcon color="cyan" variant="light" size={40} radius="xl">
              <Cpu size={20} />
            </ThemeIcon>
            <div>
              <Text fw={700} size="md" c="white">
                Lead Director Orchestrator Loop
              </Text>
              <Text size="xs" c="dimmed">
                Observe → Assess → Plan → Delegate → Evaluate → Decide
              </Text>
            </div>
          </Group>

          <Badge color="cyan" variant="dot" size="md">
            {status}
          </Badge>
        </Group>

        <Paper p="sm" radius="lg" bg="dark.8" withBorder className="border-studio-border/60">
          <Group justify="space-between" mb="xs">
            <Group gap="xs">
              <Activity size={14} className="text-cyan-400" />
              <Text size="xs" fw={700} c="cyan.4">
                Current Phase: {phase}
              </Text>
            </Group>
            <Text size="xs" c="dimmed">
              Step {stepCount} / 6
            </Text>
          </Group>
          <Progress value={(stepCount / 6) * 100} color="cyan" size="sm" radius="xl" animated />
        </Paper>
      </Stack>
    </Paper>
  );
}
