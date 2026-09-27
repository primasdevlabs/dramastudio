"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { BookOpen, Save, CheckCircle2 } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, Button, TextInput, Textarea, ThemeIcon, Alert, SimpleGrid } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { SeriesBible } from "@/lib/api/types";

export default function SeriesBiblePage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const { data: bible } = useQuery({
    queryKey: ["bible", projectId],
    queryFn: () => api.get<SeriesBible>(`/v1/projects/${projectId}/bible`).catch(() => null),
  });

  const [premise, setPremise] = useState("");
  const [genre, setGenre] = useState("Drama/Thriller");
  const [themes, setThemes] = useState("Betrayal, Power, Redemption");
  const [worldRules, setWorldRules] = useState("No magic, Grounded realism");
  const [narrativeRules, setNarrativeRules] = useState("Actions have permanent consequences");
  const [savedSuccess, setSavedSuccess] = useState(false);

  useEffect(() => {
    if (bible) {
      setPremise(bible.premise || "");
      setGenre(bible.genre || "Drama/Thriller");
      setThemes(bible.themes?.join(", ") || "");
      setWorldRules(bible.world_rules?.join("\n") || "");
      setNarrativeRules(bible.narrative_rules?.join("\n") || "");
    }
  }, [bible]);

  const saveMutation = useMutation({
    mutationFn: async () => {
      return api.post<SeriesBible>(`/v1/projects/${projectId}/bible`, {
        premise,
        genre,
        themes: themes.split(",").map((t) => t.trim()).filter(Boolean),
        world_rules: worldRules.split("\n").map((r) => r.trim()).filter(Boolean),
        narrative_rules: narrativeRules.split("\n").map((r) => r.trim()).filter(Boolean),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bible", projectId] });
      setSavedSuccess(true);
      notifications.show({
        title: "Bible Saved",
        message: "Series Bible updated to next canonical version",
        color: "amber",
      });
      setTimeout(() => setSavedSuccess(false), 3000);
    },
  });

  return (
    <Stack gap="lg" className="max-w-5xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="xl" withBorder className="bg-studio-card border-studio-border shadow-xl">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="amber" variant="light" size={44} radius="xl">
              <BookOpen size={22} />
            </ThemeIcon>
            <div>
              <Group gap="xs">
                <Title order={3} c="white">
                  Series Bible Workspace
                </Title>
                {bible && (
                  <Badge color="amber" variant="light" size="sm">
                    Version {bible.version}
                  </Badge>
                )}
              </Group>
              <Text size="xs" c="dimmed">
                Canonical creative foundation for AI generation
              </Text>
            </div>
          </Group>

          <Button
            onClick={() => saveMutation.mutate()}
            loading={saveMutation.isPending}
            leftSection={<Save size={16} />}
            variant="gradient"
            gradient={{ from: "amber", to: "orange", deg: 90 }}
            color="amber"
          >
            Save Bible Version
          </Button>
        </Group>
      </Paper>

      {savedSuccess && (
        <Alert icon={<CheckCircle2 size={16} />} title="Canonical Version Updated" color="green">
          Series Bible successfully saved and updated to next canonical version!
        </Alert>
      )}

      {/* Editor Form */}
      <Paper p="xl" radius="xl" withBorder className="bg-studio-card border-studio-border">
        <Stack gap="lg">
          <Textarea
            label="Core Story Premise"
            description="Central logline and premise for AI episode planning"
            rows={4}
            value={premise}
            onChange={(e) => setPremise(e.currentTarget.value)}
            placeholder="Write the central logline and core premise..."
            variant="filled"
          />

          <SimpleGrid cols={{ base: 1, md: 2 }} spacing="lg">
            <TextInput
              label="Genre & Tone"
              value={genre}
              onChange={(e) => setGenre(e.currentTarget.value)}
              variant="filled"
            />

            <TextInput
              label="Central Themes (Comma Separated)"
              value={themes}
              onChange={(e) => setThemes(e.currentTarget.value)}
              variant="filled"
            />
          </SimpleGrid>

          <SimpleGrid cols={{ base: 1, md: 2 }} spacing="lg">
            <Textarea
              label="World Rules (One per line)"
              rows={4}
              value={worldRules}
              onChange={(e) => setWorldRules(e.currentTarget.value)}
              placeholder="e.g. Grounded realism&#10;No supernatural events"
              variant="filled"
            />

            <Textarea
              label="Narrative Rules (One per line)"
              rows={4}
              value={narrativeRules}
              onChange={(e) => setNarrativeRules(e.currentTarget.value)}
              placeholder="e.g. Actions have permanent consequences"
              variant="filled"
            />
          </SimpleGrid>
        </Stack>
      </Paper>
    </Stack>
  );
}
