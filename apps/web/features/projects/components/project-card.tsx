import Link from "next/link";
import { Clapperboard, DollarSign, ChevronRight, Activity } from "lucide-react";
import { Card, Badge, Group, Text, Button, ThemeIcon, Stack } from "@mantine/core";
import { Project } from "@/lib/api/types";

export function ProjectCard({ project }: { project: Project }) {
  return (
    <Card
      padding="lg"
      radius="md"
      className="bg-studio-card border border-studio-border hover:border-studio-accent/60 transition-all flex flex-col justify-between"
    >
      <Stack gap="xs">
        <Group justify="space-between">
          <Group gap="xs">
            <ThemeIcon variant="light" color="terracotta" size="md" radius="sm">
              <Clapperboard size={16} />
            </ThemeIcon>
            <Badge color="terracotta" variant="outline" size="sm">
              {project.genre || "Drama"}
            </Badge>
          </Group>

          <Badge color="emerald" variant="dot" size="sm">
            {project.status || "ACTIVE"}
          </Badge>
        </Group>

        <div>
          <Text fw={700} size="lg" c="white" className="tracking-tight">
            {project.name}
          </Text>
          <Text size="xs" c="dimmed" className="line-clamp-2 mt-1">
            {project.description || "Serialized AI Drama Production"}
          </Text>
        </div>
      </Stack>

      <Group justify="space-between" align="center" className="pt-4 border-t border-studio-border/60 mt-4">
        <Group gap="md">
          <Group gap={4}>
            <Activity size={14} className="text-studio-accent" />
            <Text size="xs" c="dimmed" fw={600} className="font-mono">
              {project.mode.toUpperCase()}
            </Text>
          </Group>

          <Group gap={4}>
            <DollarSign size={14} className="text-emerald-400" />
            <Text size="xs" c="dimmed" className="font-mono">
              ${project.budget?.current_spent || 0}
            </Text>
          </Group>
        </Group>

        <Button
          component={Link}
          href={`/projects/${project.id}`}
          variant="subtle"
          color="terracotta"
          size="xs"
          rightSection={<ChevronRight size={14} />}
        >
          Open Studio
        </Button>
      </Group>
    </Card>
  );
}
