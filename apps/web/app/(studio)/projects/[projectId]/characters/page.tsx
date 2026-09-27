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
        color: "pink",
      });
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="xl" withBorder className="bg-studio-card border-studio-border shadow-xl">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="pink" variant="light" size={44} radius="xl">
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
            variant="gradient"
            gradient={{ from: "pink", to: "grape", deg: 90 }}
            color="pink"
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
        overlayProps={{ backgroundOpacity: 0.7, blur: 8 }}
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
              color="pink"
            >
              Save Character
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Grid */}
      <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="lg">
        {isLoading ? (
          [1, 2].map((i) => <Paper key={i} h={180} radius="lg" className="bg-studio-card animate-pulse" />)
        ) : characters && characters.length > 0 ? (
          characters.map((c) => (
            <Card
              key={c.id}
              p="lg"
              radius="lg"
              withBorder
              className="bg-studio-card border-studio-border hover:border-pink-500/40 transition-all flex flex-col justify-between"
            >
              <Stack gap="sm">
                <Group justify="space-between">
                  <Group gap="sm">
                    <Avatar color="pink" radius="xl" size="md">
                      {c.name.substring(0, 2).toUpperCase()}
                    </Avatar>
                    <div>
                      <Text fw={700} size="sm" c="white">
                        {c.name}
                      </Text>
                      <Text size="xs" c="pink.4" fw={600}>
                        {c.role}
                      </Text>
                    </div>
                  </Group>

                  <Badge color="pink" variant="light" size="sm">
                    v{c.version || 1}
                  </Badge>
                </Group>

                <Paper p="xs" radius="md" bg="dark.7" withBorder className="border-studio-border/60">
                  <Text size="xs" c="dimmed" className="line-clamp-3">
                    {c.bio || "No backstory notes configured."}
                  </Text>
                </Paper>
              </Stack>

              <Group justify="space-between" className="pt-3 border-t border-studio-border/60 mt-4">
                <Group gap={4}>
                  <Shirt size={14} className="text-pink-400" />
                  <Text size="xs" c="dimmed">
                    Wardrobe Locked
                  </Text>
                </Group>
                <Group gap={4}>
                  <UserCheck size={14} className="text-emerald-400" />
                  <Text size="xs" c="emerald.4" fw={600}>
                    Canon Reference
                  </Text>
                </Group>
              </Group>
            </Card>
          ))
        ) : (
          <Paper p="xl" radius="xl" withBorder className="col-span-full bg-studio-card border-studio-border text-center">
            <Text size="xs" c="dimmed">
              No characters created yet. Add a character to build your cast.
            </Text>
          </Paper>
        )}
      </SimpleGrid>
    </Stack>
  );
}
