"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Globe, Plus, Sun, Moon, CloudRain, MapPin } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, Button, Modal, TextInput, Select, Textarea, ThemeIcon, SimpleGrid, Card } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { Location } from "@/lib/api/types";

export default function WorldStudioPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [type, setType] = useState("Interior");
  const [opened, setOpened] = useState(false);

  const { data: locations, isLoading } = useQuery({
    queryKey: ["locations", projectId],
    queryFn: async () => {
      const res = await api.get<{ locations: Location[] }>(`/v1/world/locations?project_id=${projectId}`);
      return res.locations || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      return api.post<Location>("/v1/world/locations", {
        project_id: projectId,
        name,
        description,
        type,
      });
    },
    onSuccess: (loc) => {
      queryClient.invalidateQueries({ queryKey: ["locations", projectId] });
      setOpened(false);
      setName("");
      setDescription("");
      notifications.show({
        title: "Location Registered",
        message: `Registered location ${loc.name || name}`,
        color: "emerald",
      });
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="emerald" variant="light" size={44} radius="md">
              <Globe size={22} />
            </ThemeIcon>
            <div>
              <Title order={3} c="white">
                World Studio
              </Title>
              <Text size="xs" c="dimmed">
                Manage persistent production locations, environments, and lighting variants
              </Text>
            </div>
          </Group>

          <Button
            onClick={() => setOpened(true)}
            leftSection={<Plus size={16} />}
            variant="filled"
            color="emerald"
            size="sm"
            radius="sm"
          >
            Add Production Location
          </Button>
        </Group>
      </Paper>

      {/* Modal */}
      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title={
          <Text fw={700} size="md" c="white">
            New Production Location
          </Text>
        }
        centered
        overlayProps={{ backgroundOpacity: 0.7, blur: 4 }}
      >
        <Stack gap="md">
          <TextInput
            label="Location Name"
            placeholder="e.g. Downtown Apartment"
            required
            value={name}
            onChange={(e) => setName(e.currentTarget.value)}
            variant="filled"
          />

          <Select
            label="Type"
            value={type}
            onChange={(val) => setType(val || "Interior")}
            data={[
              { value: "Interior", label: "Interior" },
              { value: "Exterior", label: "Exterior" },
              { value: "Studio Set", label: "Studio Set" },
            ]}
            variant="filled"
          />

          <Textarea
            label="Description & Visual Rules"
            placeholder="Describe lighting, decor, and spatial layout..."
            rows={3}
            value={description}
            onChange={(e) => setDescription(e.currentTarget.value)}
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
              color="emerald"
            >
              Save Location
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Grid */}
      <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="lg">
        {isLoading ? (
          [1, 2].map((i) => <Paper key={i} h={170} radius="md" className="bg-studio-card animate-pulse" />)
        ) : locations && locations.length > 0 ? (
          locations.map((loc) => (
            <Card
              key={loc.id}
              p="lg"
              radius="md"
              withBorder
              className="bg-studio-card border-studio-border hover:border-emerald-500/40 transition-all flex flex-col justify-between"
            >
              <Stack gap="sm">
                <Group justify="space-between">
                  <Group gap="xs">
                    <MapPin size={16} className="text-emerald-400" />
                    <Text fw={700} size="md" c="white">
                      {loc.name}
                    </Text>
                  </Group>

                  <Badge color="emerald" variant="light" size="sm">
                    {loc.type}
                  </Badge>
                </Group>

                <Paper p="xs" radius="sm" className="bg-studio-panel border border-studio-border">
                  <Text size="xs" c="dimmed" className="line-clamp-2">
                    {loc.description || "No visual layout description specified."}
                  </Text>
                </Paper>
              </Stack>

              <Group justify="space-between" className="pt-3 border-t border-studio-border/60 mt-4">
                <Text size="xs" fw={600} c="white">
                  Lighting Variants:
                </Text>
                <Group gap="xs">
                  <Badge color="amber" variant="subtle" size="xs" leftSection={<Sun size={10} />}>
                    Day
                  </Badge>
                  <Badge color="gray" variant="subtle" size="xs" leftSection={<Moon size={10} />}>
                    Night
                  </Badge>
                  <Badge color="gray" variant="subtle" size="xs" leftSection={<CloudRain size={10} />}>
                    Rain
                  </Badge>
                </Group>
              </Group>
            </Card>
          ))
        ) : (
          <Paper p="xl" radius="md" withBorder className="col-span-full bg-studio-card border-studio-border text-center">
            <Text size="xs" c="dimmed">
              No locations registered yet. Add a location to set the scene.
            </Text>
          </Paper>
        )}
      </SimpleGrid>
    </Stack>
  );
}
