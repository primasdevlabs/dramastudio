"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Users, Plus, Shirt, UserCheck } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, Button, Modal, TextInput, Textarea, ThemeIcon, Avatar, SimpleGrid, Card } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { Character } from "@/lib/api/types";

export default function CharacterStudioPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [name, setName] = useState("");
  const [role, setRole] = useState("Protagonist");
  const [bio, setBio] = useState("");
  const [opened, setOpened] = useState(false);

  const { data: characters, isLoading } = useQuery({
    queryKey: ["characters", projectId],
    queryFn: async () => {
      const res = await api.get<{ characters: Character[] }>("/v1/characters");
      return res.characters || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      return api.post<Character>("/v1/characters", {
        project_id: projectId,
        name,
        role,
        bio,
      });
    },
    onSuccess: (c) => {
      queryClient.invalidateQueries({ queryKey: ["characters", projectId] });
      setOpened(false);
      setName("");
      setBio("");
      notifications.show({
        title: "Character Created",
        message: `Added character profile for ${c.name || name}`,
        color: "terracotta",
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
              <Users size={22} />
            </ThemeIcon>
            <div>
              <Title order={3} c="white">
                Character Studio
              </Title>
              <Text size="xs" c="dimmed">
                Manage persistent character identities, wardrobe, and voice profiles
              </Text>
            </div>
          </Group>

          <Button
            onClick={() => setOpened(true)}
            leftSection={<Plus size={16} />}
            variant="filled"
            color="terracotta"
            size="sm"
            radius="sm"
          >
            Add Character Profile
          </Button>
        </Group>
      </Paper>

      {/* Modal */}
      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title={
          <Text fw={700} size="md" c="white">
            New Character Profile
          </Text>
        }
        centered
        overlayProps={{ backgroundOpacity: 0.7, blur: 4 }}
      >
        <Stack gap="md">
          <TextInput
            label="Character Name"
            placeholder="e.g. Sarah Johnson"
            required
            value={name}
            onChange={(e) => setName(e.currentTarget.value)}
            variant="filled"
          />

          <TextInput
            label="Role / Archetype"
            placeholder="e.g. Protagonist / Journalist"
            value={role}
            onChange={(e) => setRole(e.currentTarget.value)}
            variant="filled"
          />

          <Textarea
            label="Biography & Appearance Notes"
            placeholder="Character backstory and visual constraints..."
            rows={3}
            value={bio}
            onChange={(e) => setBio(e.currentTarget.value)}
            variant="filled"
          />

          <Group justify="flex-end" gap="sm" mt="md">
            <Button variant="subtle" color="gray" onClick={() => setOpened(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => createMutation.mutate()}
              loading={createMutation.isPending}
              disabled={!name}
              color="terracotta"
            >
              Save Character
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Grid */}
      <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="lg">
        {isLoading ? (
          [1, 2].map((i) => <Paper key={i} h={180} radius="md" className="bg-studio-card animate-pulse" />)
        ) : characters && characters.length > 0 ? (
          characters.map((c) => (
            <Card
              key={c.id}
              p="lg"
              radius="md"
              withBorder
              className="bg-studio-card border-studio-border hover:border-studio-accent/60 transition-all flex flex-col justify-between"
            >
              <Stack gap="sm">
                <Group justify="space-between">
                  <Group gap="sm">
                    <Avatar color="terracotta" radius="md" size="md">
                      {c.name.substring(0, 2).toUpperCase()}
                    </Avatar>
                    <div>
                      <Text fw={700} size="sm" c="white">
                        {c.name}
                      </Text>
                      <Text size="xs" c="terracotta.4" fw={600}>
                        {c.role}
                      </Text>
                    </div>
                  </Group>

                  <Badge color="terracotta" variant="light" size="sm" className="font-mono">
                    v{c.version || 1}
                  </Badge>
                </Group>

                <Paper p="xs" radius="sm" className="bg-studio-panel border border-studio-border">
                  <Text size="xs" c="dimmed" className="line-clamp-3">
                    {c.bio || "No backstory notes configured."}
                  </Text>
                </Paper>
              </Stack>

              <Group justify="space-between" className="pt-3 border-t border-studio-border/60 mt-4">
                <Group gap={4}>
                  <Shirt size={14} className="text-amber-400" />
                  <Text size="xs" c="dimmed">
                    Wardrobe Locked
                  </Text>
                </Group>
                <Group gap={4}>
                  <UserCheck size={14} className="text-emerald-400" />
                  <Text size="xs" c="emerald.4" fw={600} className="font-mono">
                    Canon Reference
                  </Text>
                </Group>
              </Group>
            </Card>
          ))
        ) : (
          <Paper p="xl" radius="md" withBorder className="col-span-full bg-studio-card border-studio-border text-center">
            <Text size="xs" c="dimmed">
              No characters created yet. Add a character to build your cast.
            </Text>
          </Paper>
        )}
      </SimpleGrid>
    </Stack>
  );
}
