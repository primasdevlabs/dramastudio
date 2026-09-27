"use client";

import { Paper, Group, Stack, Text, Badge, ThemeIcon, Progress } from "@mantine/core";
import { Cpu, Activity } from "lucide-react";

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
    <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
      <Stack gap="md">
        <Group justify="space-between">
          <Group gap="md">
            <ThemeIcon color="terracotta" variant="light" size={40} radius="sm">
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

          <Badge color="terracotta" variant="dot" size="md">
            {status}
          </Badge>
        </Group>

        <Paper p="sm" radius="sm" className="bg-studio-panel border border-studio-border">
          <Group justify="space-between" mb="xs">
            <Group gap="xs">
              <Activity size={14} className="text-studio-accent" />
              <Text size="xs" fw={700} c="terracotta.4" className="font-mono">
                Current Phase: {phase}
              </Text>
            </Group>
            <Text size="xs" c="dimmed" className="font-mono">
              Step {stepCount} / 6
            </Text>
          </Group>
          <Progress value={(stepCount / 6) * 100} color="terracotta" size="sm" radius="xs" />
        </Paper>
      </Stack>
    </Paper>
  );
}
