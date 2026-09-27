"use client";

import { useState } from "react";
import { Plus, Clapperboard } from "lucide-react";
import { Modal, Button, TextInput, Textarea, Select, Group, Stack, Text, ThemeIcon } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { Project } from "@/lib/api/types";

export function NewProjectDialog({ onCreated }: { onCreated?: (p: Project) => void }) {
  const [opened, setOpened] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [genre, setGenre] = useState("Drama/Thriller");
  const [language, setLanguage] = useState("en");
  const [mode, setMode] = useState<"monitored" | "autonomous">("monitored");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setIsSubmitting(true);
    try {
      const proj = await api.post<Project>("/v1/projects", {
        name,
        description,
        genre,
        language,
        mode,
      });
      setOpened(false);
      setName("");
      setDescription("");
      notifications.show({
        title: "Production Initialized",
        message: `Created series ${proj.name || name} successfully`,
        color: "terracotta",
      });
      if (onCreated) onCreated(proj);
    } catch (err) {
      console.error("Failed to create project", err);
      notifications.show({
        title: "Creation Error",
        message: "Failed to initialize production series",
        color: "red",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <Button
        onClick={() => setOpened(true)}
        leftSection={<Plus size={16} />}
        color="terracotta"
        variant="filled"
        size="sm"
        radius="sm"
      >
        New Production
      </Button>

      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title={
          <Group gap="xs">
            <ThemeIcon color="terracotta" variant="light" size="lg" radius="sm">
              <Clapperboard size={18} />
            </ThemeIcon>
            <div>
              <Text fw={700} size="md" c="white">
                Create Drama Production
              </Text>
              <Text size="xs" c="dimmed">
                Set up your new AI drama series workspace
              </Text>
            </div>
          </Group>
        }
        centered
        size="md"
        overlayProps={{ backgroundOpacity: 0.7, blur: 4 }}
      >
        <form onSubmit={handleSubmit}>
          <Stack gap="md" mt="sm">
            <TextInput
              label="Production Title"
              placeholder="e.g. The Last Promise"
              required
              value={name}
              onChange={(e) => setName(e.currentTarget.value)}
              variant="filled"
            />

            <Textarea
              label="Premise / Description"
              placeholder="Brief story premise..."
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.currentTarget.value)}
              variant="filled"
            />

            <Group grow gap="md">
              <Select
                label="Genre"
                value={genre}
                onChange={(val) => setGenre(val || "Drama/Thriller")}
                data={[
                  { value: "Drama/Thriller", label: "Drama / Thriller" },
                  { value: "Romance", label: "Romance" },
                  { value: "Sci-Fi", label: "Sci-Fi" },
                  { value: "Action", label: "Action" },
                  { value: "Mystery", label: "Mystery" },
                ]}
                variant="filled"
              />

              <Select
                label="Production Mode"
                value={mode}
                onChange={(val) => setMode((val as any) || "monitored")}
                data={[
                  { value: "monitored", label: "Monitored (Approval)" },
                  { value: "autonomous", label: "Autonomous (Director)" },
                ]}
                variant="filled"
              />
            </Group>

            <Group justify="flex-end" gap="sm" mt="lg">
              <Button variant="subtle" color="gray" onClick={() => setOpened(false)}>
                Cancel
              </Button>
              <Button type="submit" loading={isSubmitting} color="terracotta">
                Initialize Production
              </Button>
            </Group>
          </Stack>
        </form>
      </Modal>
    </>
  );
}
