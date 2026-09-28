"use client";

import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Settings,
  Cpu,
  Shield,
  Key,
  Bell,
  Globe,
  DollarSign,
  Users,
  Sliders,
  ChevronRight,
  AlertTriangle,
  Check,
  RefreshCw,
  Info,
  CheckCircle2,
  Save,
  Link2,
} from "lucide-react";
import {
  Paper,
  Group,
  Stack,
  Title,
  Text,
  Badge,
  ThemeIcon,
  Tabs,
  Switch,
  Button,
  TextInput,
  NumberInput,
  Select,
  Progress,
  Divider,
  Alert,
} from "@mantine/core";
import Link from "next/link";
import { api } from "@/lib/api/client";
import type { Project } from "@/lib/api/types";
import { useSettingsStore } from "@/stores/settings-store";

interface CapabilityMapping {
  capability: string;
  model: string;
  provider: string;
  status: "active" | "degraded" | "offline";
}

interface TeamMember {
  name: string;
  email: string;
  role: string;
  lastActive: string;
}

export default function ProjectSettingsPage({
  params,
}: {
  params: { projectId: string };
}) {
  const { projectId } = params;

  const {
    globalSettings,
    getEffectiveSettings,
    getProjectOverrideState,
    toggleProjectInheritance,
    updateProjectOverride,
  } = useSettingsStore();

  const overrideState = getProjectOverrideState(projectId);
  const effectiveSettings = getEffectiveSettings(projectId);

  const [isInheriting, setIsInheriting] = useState(overrideState.inheritGlobalDefaults);
  const [projectOverrides, setProjectOverrides] = useState(overrideState.overrides);
  const [saveSuccess, setSaveSuccess] = useState(false);

  useEffect(() => {
    setIsInheriting(overrideState.inheritGlobalDefaults);
    setProjectOverrides(overrideState.overrides);
  }, [projectId, overrideState.inheritGlobalDefaults, overrideState.overrides]);

  const { data: project } = useQuery({
    queryKey: ["project", projectId],
    queryFn: () => api.get<Project>(`/v1/projects/${projectId}`).catch(() => null),
  });

  const { data: capabilities } = useQuery({
    queryKey: ["capabilities", projectId],
    queryFn: async () => {
      const res = await api
        .get<{ capabilities: CapabilityMapping[] }>(
          `/v1/projects/${projectId}/settings/capabilities`
        )
        .catch(() => null);

      return (
        res?.capabilities ?? [
          { capability: "Scriptwriting & Dialogue", model: effectiveSettings.scriptModel, provider: "OpenAI", status: "active" as const },
          { capability: "Video Shot Generation", model: effectiveSettings.videoGenModel, provider: "Runway", status: "active" as const },
          { capability: "Voice & Dialogue Synthesis", model: effectiveSettings.characterVoiceEngine, provider: "ElevenLabs", status: "active" as const },
          { capability: "Storyboard Layout", model: "StoryboardAgent-v1", provider: "Internal", status: "active" as const },
          { capability: "Continuity Verification", model: "ContinuityChecker-v3", provider: "Internal", status: "active" as const },
          { capability: "Image Generation", model: effectiveSettings.imageGenModel, provider: "Midjourney", status: "active" as const },
        ]
      );
    },
  });

  const handleToggleInherit = (inherit: boolean) => {
    setIsInheriting(inherit);
    toggleProjectInheritance(projectId, inherit);
  };

  const handleSave = () => {
    updateProjectOverride(projectId, projectOverrides);
    setSaveSuccess(true);
    setTimeout(() => setSaveSuccess(false), 3000);
  };

  const statusColors: Record<string, string> = {
    active: "emerald",
    degraded: "amber",
    offline: "gray",
  };

  const budget = project?.budget ?? {
    total_budget: effectiveSettings.totalProjectBudgetCap,
    current_spent: 127.45,
    max_cost_per_generation: effectiveSettings.maxCostPerGeneration,
    currency: effectiveSettings.currency,
  };

  const budgetPercent = Math.round((budget.current_spent / budget.total_budget) * 100);

  const teamMembers: TeamMember[] = [
    { name: "Studio Operator", email: "operator@dramastudio.ai", role: "Owner", lastActive: "Now" },
    { name: "Lead Director", email: "director@system", role: "Director", lastActive: "Active" },
    { name: "Story Director", email: "story@system", role: "Director", lastActive: "Active" },
  ];

  return (
    <Stack gap="lg" className="max-w-5xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="md" withBorder className="bg-studio-card border-studio-border">
        <Group justify="space-between" align="flex-start">
          <Group gap="md">
            <ThemeIcon color="terracotta" variant="light" size={44} radius="md">
              <Settings size={22} />
            </ThemeIcon>
            <div>
              <div className="flex items-center gap-2">
                <Title order={3} c="white">
                  Project Settings & Configuration
                </Title>
                <Badge
                  color={isInheriting ? "blue" : "terracotta"}
                  variant="light"
                  size="sm"
                >
                  {isInheriting ? "INHERITING GLOBAL DEFAULTS" : "PROJECT OVERRIDES ACTIVE"}
                </Badge>
              </div>
              <Text size="xs" c="dimmed">
                Configure production rules, model selection, aspect ratio, and budget limits for this specific project.
              </Text>
            </div>
          </Group>

          <Button
            color="terracotta"
            size="xs"
            leftSection={<Save size={14} />}
            onClick={handleSave}
            disabled={isInheriting}
          >
            Save Project Overrides
          </Button>
        </Group>
      </Paper>

      {saveSuccess && (
        <Alert
          icon={<CheckCircle2 size={16} />}
          title="Project Settings Saved"
          color="green"
          variant="light"
          withCloseButton
          onClose={() => setSaveSuccess(false)}
        >
          Project-specific configuration saved successfully.
        </Alert>
      )}

      {/* Global Inheritance Banner */}
      <Paper p="md" radius="md" className="bg-studio-panel border border-studio-border">
        <Group justify="space-between" align="center">
          <Group gap="sm">
            <Info size={18} className="text-studio-accent shrink-0" />
            <div>
              <Text fw={600} size="xs" c="white">
                Global Studio Defaults Inheritance
              </Text>
              <Text size="xs" c="dimmed">
                {isInheriting
                  ? "This project is automatically adopting all Global Studio Settings. Disable inheritance below to customize specific settings."
                  : "Inheritance disabled. Custom project-level overrides are active for this production."}
              </Text>
            </div>
          </Group>

          <Group gap="md">
            <Link href="/settings">
              <Button variant="subtle" color="gray" size="xs" leftSection={<Link2 size={12} />}>
                View Global Settings
              </Button>
            </Link>
            <Switch
              color="terracotta"
              label="Inherit Global Defaults"
              checked={isInheriting}
              onChange={(e) => handleToggleInherit(e.currentTarget.checked)}
            />
          </Group>
        </Group>
      </Paper>

      <Tabs defaultValue="models" variant="outline">
        <Tabs.List className="border-b border-studio-border">
          <Tabs.Tab value="models" leftSection={<Cpu size={14} />}>
            Production Models & Engines
          </Tabs.Tab>
          <Tabs.Tab value="production" leftSection={<Sliders size={14} />}>
            Production & Output
          </Tabs.Tab>
          <Tabs.Tab value="capabilities" leftSection={<Cpu size={14} />}>
            Capability Providers
          </Tabs.Tab>
          <Tabs.Tab value="budgets" leftSection={<DollarSign size={14} />}>
            Budget & Costs
          </Tabs.Tab>
          <Tabs.Tab value="team" leftSection={<Users size={14} />}>
            Team & Permissions
          </Tabs.Tab>
          <Tabs.Tab value="notifications" leftSection={<Bell size={14} />}>
            Notifications
          </Tabs.Tab>
        </Tabs.List>

        {/* Models Tab */}
        <Tabs.Panel value="models" pt="lg">
          <Stack gap="md">
            <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border space-y-4">
              <div className="flex justify-between items-center">
                <Text fw={700} size="sm" c="white">
                  Model Selection Overrides
                </Text>
                <Badge color={isInheriting ? "gray" : "terracotta"} variant="outline" size="xs">
                  {isInheriting ? "INHERITED FROM GLOBAL" : "CUSTOM OVERRIDES"}
                </Badge>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Select
                  label="Scriptwriting & Screenplay Model"
                  value={isInheriting ? globalSettings.scriptModel : (projectOverrides.scriptModel || globalSettings.scriptModel)}
                  disabled={isInheriting}
                  onChange={(val) => val && setProjectOverrides((prev) => ({ ...prev, scriptModel: val }))}
                  data={[
                    { value: "gpt-4o", label: "OpenAI GPT-4o" },
                    { value: "claude-3-5-sonnet", label: "Anthropic Claude 3.5 Sonnet" },
                    { value: "gemini-1.5-pro", label: "Google Gemini 1.5 Pro" },
                  ]}
                />

                <Select
                  label="World Building Intelligence"
                  value={isInheriting ? globalSettings.worldBuildingModel : (projectOverrides.worldBuildingModel || globalSettings.worldBuildingModel)}
                  disabled={isInheriting}
                  onChange={(val) => val && setProjectOverrides((prev) => ({ ...prev, worldBuildingModel: val }))}
                  data={[
                    { value: "claude-3-5-sonnet", label: "Anthropic Claude 3.5 Sonnet" },
                    { value: "gpt-4o", label: "OpenAI GPT-4o" },
                    { value: "gemini-1.5-pro", label: "Google Gemini 1.5 Pro" },
                  ]}
                />

                <Select
                  label="Character Voice Engine"
                  value={isInheriting ? globalSettings.characterVoiceEngine : (projectOverrides.characterVoiceEngine || globalSettings.characterVoiceEngine)}
                  disabled={isInheriting}
                  onChange={(val) => val && setProjectOverrides((prev) => ({ ...prev, characterVoiceEngine: val }))}
                  data={[
                    { value: "elevenlabs-multilingual-v2", label: "ElevenLabs Multilingual v2" },
                    { value: "openai-tts-1-hd", label: "OpenAI TTS-1 HD" },
                    { value: "bark-v2", label: "Suno Bark Open-Source" },
                  ]}
                />

                <Select
                  label="Video Generation Model"
                  value={isInheriting ? globalSettings.videoGenModel : (projectOverrides.videoGenModel || globalSettings.videoGenModel)}
                  disabled={isInheriting}
                  onChange={(val) => val && setProjectOverrides((prev) => ({ ...prev, videoGenModel: val }))}
                  data={[
                    { value: "runway-gen3-alpha", label: "Runway Gen-3 Alpha" },
                    { value: "luma-dream-machine", label: "Luma Dream Machine" },
                    { value: "pika-2", label: "Pika 2.0" },
                  ]}
                />
              </div>
            </Paper>
          </Stack>
        </Tabs.Panel>

        {/* Production Tab */}
        <Tabs.Panel value="production" pt="lg">
          <Stack gap="md">
            <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
              <Text fw={700} size="sm" c="white" mb="md">
                Production Mode
              </Text>
              <Stack gap="md">
                <Group justify="space-between">
                  <div>
                    <Text fw={600} size="xs" c="white">Autonomy Mode</Text>
                    <Text size="xs" c="dimmed">
                      Switch between monitored (human-in-the-loop) and autonomous production
                    </Text>
                  </div>
                  <Select
                    data={[
                      { value: "monitored", label: "Monitored" },
                      { value: "autonomous", label: "Autonomous" },
                    ]}
                    value={project?.mode ?? "monitored"}
                    variant="filled"
                    size="xs"
                    w={160}
                  />
                </Group>
                <Divider color="dark.5" />
                <Group justify="space-between">
                  <div>
                    <Text fw={600} size="xs" c="white">Human-in-the-Loop Approval</Text>
                    <Text size="xs" c="dimmed">
                      Require human approval before advancing production stages
                    </Text>
                  </div>
                  <Switch color="terracotta" defaultChecked />
                </Group>
                <Group justify="space-between">
                  <div>
                    <Text fw={600} size="xs" c="white">Automatic Continuity Checks</Text>
                    <Text size="xs" c="dimmed">
                      Run continuity verification after each generation
                    </Text>
                  </div>
                  <Switch color="terracotta" defaultChecked />
                </Group>
              </Stack>
            </Paper>

            <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
              <Text fw={700} size="sm" c="white" mb="md">
                Output Configuration
              </Text>
              <Stack gap="md">
                <Group gap="md" grow>
                  <Select
                    label="Default Aspect Ratio"
                    data={["9:16", "16:9", "1:1", "4:3"]}
                    value={isInheriting ? globalSettings.defaultAspectRatio : (projectOverrides.defaultAspectRatio || globalSettings.defaultAspectRatio)}
                    disabled={isInheriting}
                    onChange={(val) => val && setProjectOverrides((prev) => ({ ...prev, defaultAspectRatio: val as any }))}
                    variant="filled"
                    size="xs"
                  />
                  <Select
                    label="Default Resolution"
                    data={["720p", "1080p", "4K"]}
                    value={isInheriting ? globalSettings.defaultResolution : (projectOverrides.defaultResolution || globalSettings.defaultResolution)}
                    disabled={isInheriting}
                    onChange={(val) => val && setProjectOverrides((prev) => ({ ...prev, defaultResolution: val as any }))}
                    variant="filled"
                    size="xs"
                  />
                </Group>
              </Stack>
            </Paper>
          </Stack>
        </Tabs.Panel>

        {/* AI Providers Tab */}
        <Tabs.Panel value="capabilities" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white" className="font-mono uppercase tracking-wider">
                Capability Provider Mapping
              </Text>
              <Button variant="subtle" color="terracotta" size="xs" leftSection={<RefreshCw size={12} />}>
                Refresh Status
              </Button>
            </Group>
            <Stack gap="sm">
              {capabilities?.map((cap) => (
                <Paper
                  key={cap.capability}
                  p="md"
                  radius="sm"
                  className="bg-studio-panel border border-studio-border"
                >
                  <Group justify="space-between" align="center">
                    <div className="flex-1 min-w-0">
                      <Text fw={700} size="xs" c="white">
                        {cap.capability}
                      </Text>
                    </div>
                    <Group gap="sm">
                      <Badge
                        color={statusColors[cap.status] ?? "gray"}
                        variant="dot"
                        size="xs"
                        className="font-mono"
                      >
                        {cap.status}
                      </Badge>
                      <Badge color="terracotta" variant="light" size="sm" className="font-mono">
                        {cap.model}
                      </Badge>
                      <Badge color="gray" variant="outline" size="xs" className="font-mono">
                        {cap.provider}
                      </Badge>
                    </Group>
                  </Group>
                </Paper>
              ))}
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Budgets Tab */}
        <Tabs.Panel value="budgets" pt="lg">
          <Stack gap="md">
            <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
              <Text fw={700} size="sm" c="white" mb="md">
                Project Budget & Cost Guards
              </Text>
              <Stack gap="md">
                <div>
                  <Group justify="space-between" mb={4}>
                    <Text size="xs" c="dimmed">Spent</Text>
                    <Text size="xs" fw={600} c="white" className="font-mono">
                      ${budget.current_spent.toFixed(2)} / ${budget.total_budget.toFixed(2)} {budget.currency}
                    </Text>
                  </Group>
                  <Progress
                    value={budgetPercent}
                    color={budgetPercent > 80 ? "red" : budgetPercent > 60 ? "amber" : "terracotta"}
                    size="sm"
                    radius="sm"
                  />
                </div>

                <Divider color="dark.5" />

                <Group gap="md" grow>
                  <NumberInput
                    label="Project Total Budget"
                    value={budget.total_budget}
                    prefix="$"
                    variant="filled"
                    size="xs"
                    min={0}
                  />
                  <NumberInput
                    label="Max Cost Per Generation"
                    value={budget.max_cost_per_generation}
                    prefix="$"
                    variant="filled"
                    size="xs"
                    min={0}
                    step={0.1}
                    decimalScale={2}
                  />
                </Group>
              </Stack>
            </Paper>
          </Stack>
        </Tabs.Panel>

        {/* Team Tab */}
        <Tabs.Panel value="team" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Group justify="space-between" mb="md">
              <Text fw={700} size="sm" c="white">
                Team & Agents
              </Text>
              <Button variant="outline" color="terracotta" size="xs" leftSection={<Users size={12} />}>
                Invite Member
              </Button>
            </Group>
            <Stack gap="sm">
              {teamMembers.map((member) => (
                <Paper
                  key={member.email}
                  p="md"
                  radius="sm"
                  className="bg-studio-panel border border-studio-border"
                >
                  <Group justify="space-between">
                    <Group gap="md">
                      <ThemeIcon
                        color={member.role === "Agent" ? "gray" : "terracotta"}
                        variant="light"
                        size="md"
                        radius="xl"
                      >
                        {member.role === "Agent" ? (
                          <Cpu size={14} />
                        ) : (
                          <Text size="xs" fw={700}>{member.name[0]}</Text>
                        )}
                      </ThemeIcon>
                      <div>
                        <Text fw={600} size="xs" c="white">{member.name}</Text>
                        <Text size="xs" c="dimmed">{member.email}</Text>
                      </div>
                    </Group>
                    <Group gap="sm">
                      <Badge
                        color={member.role === "Owner" ? "terracotta" : "gray"}
                        variant="light"
                        size="xs"
                      >
                        {member.role}
                      </Badge>
                      <Badge color="emerald" variant="dot" size="xs" className="font-mono">
                        {member.lastActive}
                      </Badge>
                    </Group>
                  </Group>
                </Paper>
              ))}
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Notifications Tab */}
        <Tabs.Panel value="notifications" pt="lg">
          <Paper p="lg" radius="md" withBorder className="bg-studio-card border-studio-border">
            <Text fw={700} size="sm" c="white" mb="md">
              Notification Preferences
            </Text>
            <Stack gap="md">
              {[
                { label: "Generation Complete", desc: "Notify when a shot or audio generation finishes" },
                { label: "Continuity Issues", desc: "Alert on newly detected continuity violations" },
                { label: "Approval Required", desc: "Notify when a human approval gate is reached" },
              ].map((pref) => (
                <Group key={pref.label} justify="space-between">
                  <div>
                    <Text fw={600} size="xs" c="white">{pref.label}</Text>
                    <Text size="xs" c="dimmed">{pref.desc}</Text>
                  </div>
                  <Switch color="terracotta" defaultChecked />
                </Group>
              ))}
            </Stack>
          </Paper>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
