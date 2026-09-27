"use client";

import { Paper, Group, Stack, Text, Badge, ThemeIcon } from "@mantine/core";
import { Bot, CheckCircle2, Clock } from "lucide-react";

export function AgentStatusCard({
  name = "Scriptwriter Agent",
  capability = "Dialogue Generation (GPT-4o)",
  status = "idle",
}: {
  name?: string;
  capability?: string;
  status?: "idle" | "working" | "completed";
}) {
  return (
    <Paper p="md" radius="md" withBorder className="bg-studio-panel border-studio-border">
      <Group justify="space-between">
        <Group gap="sm">
          <ThemeIcon color={status === "working" ? "terracotta" : "gray"} variant="light" size="md" radius="sm">
            <Bot size={16} />
          </ThemeIcon>
          <div>
            <Text fw={700} size="xs" c="white">
              {name}
            </Text>
            <Text size="xs" c="dimmed">
              {capability}
            </Text>
          </div>
        </Group>

        <Badge
          color={status === "working" ? "terracotta" : status === "completed" ? "emerald" : "gray"}
          variant="light"
          size="xs"
          className="font-mono"
        >
          {status.toUpperCase()}
        </Badge>
      </Group>
    </Paper>
  );
}
