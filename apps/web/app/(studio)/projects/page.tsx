"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { FolderPlus, AlertCircle, Plus } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, SimpleGrid, Skeleton, Alert, ThemeIcon, Button } from "@mantine/core";
import { api } from "@/lib/api/client";
import { Project } from "@/lib/api/types";
import { ProjectCard } from "@/features/projects/components/project-card";

export default function ProjectsPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["projects"],
    queryFn: async () => {
      const res = await api.get<{ items: Project[] }>("/v1/projects");
      return res.items || [];
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header Banner */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="center">
          <Stack gap={4}>
            <Group gap="xs">
              <Badge color="terracotta" variant="light" size="sm" className="font-mono">
                Control Tower
              </Badge>
            </Group>
            <Title order={2} c="white" className="tracking-tight">
              Active Drama Productions
            </Title>
            <Text size="xs" c="dimmed">
              Manage serialized drama productions, series bibles, and release schedules
            </Text>
          </Stack>
          <Button
            component={Link}
            href="/projects/new"
            leftSection={<Plus size={16} />}
            color="terracotta"
            variant="filled"
            size="sm"
            radius="sm"
          >
            Create Series
          </Button>
        </Group>
      </Paper>

      {/* Grid List */}
      {isLoading ? (
        <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="lg">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} height={180} radius="md" />
          ))}
        </SimpleGrid>
      ) : error ? (
        <Alert icon={<AlertCircle size={16} />} title="Loading Error" color="red" variant="filled">
          Failed to load productions: {(error as Error).message}
        </Alert>
      ) : data && data.length > 0 ? (
        <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="lg">
          {data.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </SimpleGrid>
      ) : (
        <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border text-center max-w-md mx-auto">
          <Stack align="center" gap="md">
            <ThemeIcon variant="light" color="terracotta" size={48} radius="md">
              <FolderPlus size={24} />
            </ThemeIcon>
            <div>
              <Text fw={700} size="lg" c="white">
                No Productions Yet
              </Text>
              <Text size="xs" c="dimmed" mt={4}>
                Define the series direction and initialize the first production.
              </Text>
            </div>
            <Button
              component={Link}
              href="/projects/new"
              leftSection={<Plus size={16} />}
              color="terracotta"
              variant="filled"
              size="sm"
              radius="sm"
            >
              Create Series
            </Button>
          </Stack>
        </Paper>
      )}
    </Stack>
  );
}
