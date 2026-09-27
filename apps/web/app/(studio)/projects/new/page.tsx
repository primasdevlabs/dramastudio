"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Clapperboard, ArrowLeft, Layers, Video, ShieldCheck, Sparkles, Check } from "lucide-react";
import {
  Paper,
  Group,
  Stack,
  Title,
  Text,
  Badge,
  TextInput,
  Textarea,
  Select,
  Button,
  ThemeIcon,
  SimpleGrid,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { Project } from "@/lib/api/types";

export default function NewProjectPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [genre, setGenre] = useState("Drama/Thriller");
  const [language, setLanguage] = useState("en");
  const [mode, setMode] = useState<"monitored" | "autonomous">("monitored");
  const [format, setFormat] = useState("9:16 Vertical");
  const [episodesCount, setEpisodesCount] = useState("12");
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

      notifications.show({
        title: "Production Workspace Created",
        message: `Initialized ${proj.name || name} successfully`,
        color: "terracotta",
      });

      router.push(`/projects/${proj.id}`);
    } catch (err) {
      console.error("Failed to create production project", err);
      notifications.show({
        title: "Initialization Error",
        message: "Failed to initialize production series workspace",
        color: "red",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Stack gap="lg" className="max-w-4xl mx-auto py-4">
      {/* Back Button & Header Banner */}
      <Group justify="space-between" align="center">
        <Button
          component={Link}
          href="/projects"
          variant="subtle"
          color="gray"
          size="xs"
          leftSection={<ArrowLeft size={14} />}
        >
          Back to Productions
        </Button>
        <Badge color="terracotta" variant="light" size="sm" className="font-mono">
          Series Workspace Setup
        </Badge>
      </Group>

      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group gap="md">
          <ThemeIcon color="terracotta" variant="light" size={48} radius="md">
            <Clapperboard size={24} />
          </ThemeIcon>
          <div>
            <Title order={2} c="white" className="tracking-tight">
              Create Drama Production
            </Title>
            <Text size="xs" c="dimmed" mt={2}>
              Set up your new AI drama series workspace, configure model capability policies, and define story canon rules.
            </Text>
          </div>
        </Group>
      </Paper>

      {/* Main Creation Form */}
      <form onSubmit={handleSubmit}>
        <Stack gap="lg">
          {/* Section 1: Series Identity */}
          <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border space-y-4">
            <Group justify="space-between" className="border-b border-studio-border pb-3">
              <Group gap="xs">
                <Layers size={16} className="text-studio-accent" />
                <Text fw={700} size="sm" c="white" uppercase className="tracking-wider font-mono">
                  Series Identity & Premise
                </Text>
              </Group>
              <Text size="xs" c="dimmed">
                Step 1 of 2
              </Text>
            </Group>

            <Stack gap="md">
              <TextInput
                label="Production Title"
                placeholder="e.g. Echoes of Tomorrow / The Midnight Promise"
                required
                value={name}
                onChange={(e) => setName(e.currentTarget.value)}
                variant="filled"
                size="sm"
              />

              <Textarea
                label="Premise & Core Narrative Logline"
                placeholder="Describe the central story premise, protagonist goals, central conflict, and thematic focus..."
                rows={4}
                value={description}
                onChange={(e) => setDescription(e.currentTarget.value)}
                variant="filled"
                size="sm"
              />

              <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
                <Select
                  label="Genre & Tone"
                  value={genre}
                  onChange={(val) => setGenre(val || "Drama/Thriller")}
                  data={[
                    { value: "Drama/Thriller", label: "Drama / Thriller" },
                    { value: "Romance", label: "Romance & Melodrama" },
                    { value: "Sci-Fi", label: "Sci-Fi & Cyberpunk" },
                    { value: "Action", label: "Action & Suspense" },
                    { value: "Mystery", label: "Mystery & Crime" },
                  ]}
                  variant="filled"
                  size="sm"
                />

                <Select
                  label="Primary Dialogue Language"
                  value={language}
                  onChange={(val) => setLanguage(val || "en")}
                  data={[
                    { value: "en", label: "English (US/UK)" },
                    { value: "es", label: "Spanish" },
                    { value: "zh", label: "Mandarin Chinese" },
                    { value: "ja", label: "Japanese" },
                    { value: "ko", label: "Korean" },
                  ]}
                  variant="filled"
                  size="sm"
                />
              </SimpleGrid>
            </Stack>
          </Paper>

          {/* Section 2: Production & Model Policy Configuration */}
          <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border space-y-4">
            <Group justify="space-between" className="border-b border-studio-border pb-3">
              <Group gap="xs">
                <Video size={16} className="text-studio-accent" />
                <Text fw={700} size="sm" c="white" uppercase className="tracking-wider font-mono">
                  Orchestration & Media Format Specs
                </Text>
              </Group>
              <Text size="xs" c="dimmed">
                Step 2 of 2
              </Text>
            </Group>

            <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="md">
              <Select
                label="Production Mode"
                value={mode}
                onChange={(val) => setMode((val as any) || "monitored")}
                data={[
                  { value: "monitored", label: "Monitored Mode (Human Approval)" },
                  { value: "autonomous", label: "Autonomous Mode (Lead Director)" },
                ]}
                variant="filled"
                size="sm"
              />

              <Select
                label="Render Format / Viewport"
                value={format}
                onChange={(val) => setFormat(val || "9:16 Vertical")}
                data={[
                  { value: "9:16 Vertical", label: "9:16 Vertical (TikTok/Reels/Shorts)" },
                  { value: "16:9 Landscape", label: "16:9 Landscape (TV/Widescreen)" },
                ]}
                variant="filled"
                size="sm"
              />

              <Select
                label="Target Season Length"
                value={episodesCount}
                onChange={(val) => setEpisodesCount(val || "12")}
                data={[
                  { value: "6", label: "6 Episodes (Mini-series)" },
                  { value: "12", label: "12 Episodes (Standard Season)" },
                  { value: "24", label: "24 Episodes (Full Season)" },
                  { value: "150", label: "150 Episodes (Long-form Drama)" },
                ]}
                variant="filled"
                size="sm"
              />
            </SimpleGrid>

            <Paper p="md" radius="sm" className="bg-studio-panel border border-studio-border mt-2">
              <Group justify="space-between">
                <div>
                  <Text size="xs" fw={700} c="white">
                    Lead Director Loop Policy
                  </Text>
                  <Text size="xs" c="dimmed" mt={1}>
                    {mode === "monitored"
                      ? "Human approval gates enabled after script generation and video shot generation before FFmpeg assembly."
                      : "Autonomous Director Loop will execute scriptwriting, Wan video generation, voice synthesis, and FFmpeg assembly automatically."}
                  </Text>
                </div>
                <Badge color={mode === "monitored" ? "amber" : "emerald"} variant="light" size="sm" className="font-mono">
                  {mode.toUpperCase()}
                </Badge>
              </Group>
            </Paper>
          </Paper>

          {/* Action Footer Bar */}
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" align="center">
              <Button
                component={Link}
                href="/projects"
                variant="subtle"
                color="gray"
                size="sm"
              >
                Cancel
              </Button>

              <Button
                type="submit"
                loading={isSubmitting}
                disabled={!name.trim()}
                color="terracotta"
                variant="filled"
                size="sm"
                radius="sm"
                leftSection={<Check size={16} />}
              >
                Initialize Drama Production
              </Button>
            </Group>
          </Paper>
        </Stack>
      </form>
    </Stack>
  );
}
