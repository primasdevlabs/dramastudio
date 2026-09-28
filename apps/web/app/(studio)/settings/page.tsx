"use client";

import { useState } from "react";
import {
  Title,
  Text,
  Paper,
  Tabs,
  TextInput,
  Select,
  NumberInput,
  Switch,
  Button,
  Group,
  Stack,
  Badge,
  Alert,
  Divider,
  PasswordInput,
  ThemeIcon,
  Modal,
  CheckIcon,
  Loader,
} from "@mantine/core";
import {
  Sliders,
  Cpu,
  Key,
  ShieldCheck,
  RotateCcw,
  Save,
  CheckCircle2,
  Info,
  Globe,
  Settings2,
  Share2,
  Video,
  Play,
  Check,
  Link2,
  Lock,
  Building,
  HardDrive,
  Users,
  DollarSign,
  AlertTriangle,
  Radio,
  ExternalLink,
  Shield,
  Trash2,
  RefreshCw,
  Zap,
} from "lucide-react";
import { useSettingsStore } from "@/stores/settings-store";
import type { ChannelConnectionDetail, ProviderConnectionStatus } from "@/lib/api/settings-types";

export default function GlobalSettingsPage() {
  const { globalSettings, updateGlobalSettings, resetGlobalSettings } =
    useSettingsStore();

  const [formData, setFormData] = useState(globalSettings);
  const [savedSuccess, setSavedSuccess] = useState(false);

  // Manage channel modal state
  const [managingChannel, setManagingChannel] = useState<ChannelConnectionDetail | null>(null);
  const [isManageModalOpen, setIsManageModalOpen] = useState(false);

  // OAuth flow modal state
  const [connectingProvider, setConnectingProvider] = useState<"facebook" | "instagram" | "tiktok" | "youtube" | null>(null);
  const [isOAuthModalOpen, setIsOAuthModalOpen] = useState(false);

  // Testing provider connection state
  const [testingProvider, setTestingProvider] = useState<string | null>(null);
  const [testResult, setTestResult] = useState<{ provider: string; success: boolean; message: string } | null>(null);

  const [providerStatuses, setProviderStatuses] = useState<Record<string, ProviderConnectionStatus>>(
    formData.providerStatuses || {
      openai: { status: "connected", lastTested: "Just now" },
      anthropic: { status: "connected", lastTested: "Just now" },
      gemini: { status: "connected", lastTested: "Just now" },
      elevenlabs: { status: "connected", lastTested: "Just now" },
      runway: { status: "connected", lastTested: "Just now" },
      replicate: { status: "connected", lastTested: "Just now" },
      midjourney: { status: "connected", lastTested: "Just now" },
    }
  );

  const handleChange = <K extends keyof typeof globalSettings>(
    key: K,
    value: (typeof globalSettings)[K]
  ) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = () => {
    updateGlobalSettings({
      ...formData,
      providerStatuses,
    });
    setSavedSuccess(true);
    setTimeout(() => setSavedSuccess(false), 3000);
  };

  const handleReset = () => {
    resetGlobalSettings();
    setFormData(useSettingsStore.getState().globalSettings);
  };

  const handleTestProviderConnection = (providerKey: string, providerName: string, keyVal: string) => {
    setTestingProvider(providerKey);
    setTestResult(null);

    // Simulate API connection test
    setTimeout(() => {
      setTestingProvider(null);
      if (!keyVal || keyVal.trim() === "") {
        setProviderStatuses((prev) => ({
          ...prev,
          [providerKey]: { status: "disconnected", lastTested: "Just now", errorMessage: "No API key configured." },
        }));
        setTestResult({
          provider: providerName,
          success: false,
          message: `${providerName} key is missing. Please enter a valid API key.`,
        });
      } else {
        setProviderStatuses((prev) => ({
          ...prev,
          [providerKey]: { status: "connected", lastTested: "Just now" },
        }));
        setTestResult({
          provider: providerName,
          success: true,
          message: `Successfully connected to ${providerName} service endpoint.`,
        });
      }
    }, 1200);
  };

  const handleOpenManage = (channel: ChannelConnectionDetail) => {
    setManagingChannel({ ...channel });
    setIsManageModalOpen(true);
  };

  const handleSaveChannelManagement = () => {
    if (!managingChannel) return;
    const provider = managingChannel.provider;
    setFormData((prev) => ({
      ...prev,
      [provider]: managingChannel,
    }));
    updateGlobalSettings({
      [provider]: managingChannel,
    });
    setIsManageModalOpen(false);
  };

  const handleDisconnectChannel = () => {
    if (!managingChannel) return;
    const provider = managingChannel.provider;
    const updatedChannel: ChannelConnectionDetail = {
      ...managingChannel,
      connected: false,
      publishingEnabled: false,
      autoPostEnabled: false,
    };
    setFormData((prev) => ({
      ...prev,
      [provider]: updatedChannel,
    }));
    updateGlobalSettings({
      [provider]: updatedChannel,
    });
    setIsManageModalOpen(false);
  };

  const handleStartOAuthConnect = (provider: "facebook" | "instagram" | "tiktok" | "youtube") => {
    setConnectingProvider(provider);
    setIsOAuthModalOpen(true);
  };

  const handleCompleteOAuthConnect = () => {
    if (!connectingProvider) return;
    const currentChannel = formData[connectingProvider];
    const defaultAccounts: Record<string, { accountName: string; handle: string }> = {
      facebook: { accountName: "Primas Studio Page", handle: "Primas Studio Page" },
      instagram: { accountName: "Primas Studio", handle: "@studio" },
      tiktok: { accountName: "DramaStudio Creator", handle: "@dramastudio_official" },
      youtube: { accountName: "DramaStudio Official", handle: "DramaStudio Shorts" },
    };

    const updatedChannel: ChannelConnectionDetail = {
      ...currentChannel,
      connected: true,
      publishingEnabled: true,
      autoPostEnabled: true,
      accountName: defaultAccounts[connectingProvider].accountName,
      handle: defaultAccounts[connectingProvider].handle,
    };

    setFormData((prev) => ({
      ...prev,
      [connectingProvider]: updatedChannel,
    }));
    updateGlobalSettings({
      [connectingProvider]: updatedChannel,
    });

    setIsOAuthModalOpen(false);
  };

  const channelList: ChannelConnectionDetail[] = [
    formData.facebook,
    formData.instagram,
    formData.tiktok,
    formData.youtube,
  ];

  const connectedChannels = channelList.filter((c) => c.connected);
  const availableChannels = channelList.filter((c) => !c.connected);

  const providerLabels: Record<string, { title: string; desc: string; iconColor: string }> = {
    facebook: {
      title: "Facebook",
      desc: "Publish episodes and promotional clips.",
      iconColor: "blue",
    },
    instagram: {
      title: "Instagram",
      desc: "Publish episodes, reels and promotional clips.",
      iconColor: "grape",
    },
    tiktok: {
      title: "TikTok",
      desc: "Publish vertical episodes and clips.",
      iconColor: "gray",
    },
    youtube: {
      title: "YouTube",
      desc: "Publish episodes, Shorts and trailers.",
      iconColor: "red",
    },
  };

  return (
    <div className="max-w-6xl mx-auto space-y-6 pb-12">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-studio-border">
        <div>
          <div className="flex items-center gap-2">
            <Title order={2} className="text-xl font-bold text-white tracking-tight">
              Global Studio Settings
            </Title>
            <Badge color="terracotta" variant="filled" size="sm">
              STUDIO DEFAULTS
            </Badge>
          </div>
          <Text size="xs" c="dimmed" mt={2}>
            Configure studio-wide distribution connections, model provider keys, output standards, and master policies.
          </Text>
        </div>

        <Group gap="sm">
          <Button
            variant="outline"
            color="gray"
            size="xs"
            leftSection={<RotateCcw size={14} />}
            onClick={handleReset}
          >
            Reset Defaults
          </Button>
          <Button
            color="terracotta"
            size="xs"
            leftSection={<Save size={14} />}
            onClick={handleSave}
          >
            Save Global Settings
          </Button>
        </Group>
      </div>

      {savedSuccess && (
        <Alert
          icon={<CheckCircle2 size={16} />}
          title="Global Settings Saved"
          color="green"
          variant="light"
          withCloseButton
          onClose={() => setSavedSuccess(false)}
        >
          Studio-wide defaults updated. All projects set to inherit global defaults will adopt these parameters.
        </Alert>
      )}

      {/* Main Settings Navigation Tabs */}
      <Tabs defaultValue="models" variant="outline" classNames={{ list: "border-studio-border" }}>
        <Tabs.List className="mb-6 flex-wrap">
          <Tabs.Tab
            value="general"
            leftSection={<Building size={14} className="text-studio-muted" />}
          >
            General
          </Tabs.Tab>
          <Tabs.Tab
            value="production"
            leftSection={<Sliders size={14} className="text-studio-muted" />}
          >
            Production
          </Tabs.Tab>
          <Tabs.Tab
            value="models"
            leftSection={<Cpu size={14} className="text-studio-accent" />}
            className="font-bold text-white"
          >
            Models & Providers
          </Tabs.Tab>
          <Tabs.Tab
            value="storage"
            leftSection={<HardDrive size={14} className="text-studio-muted" />}
          >
            Storage
          </Tabs.Tab>
          <Tabs.Tab
            value="team"
            leftSection={<Users size={14} className="text-studio-muted" />}
          >
            Team & Access
          </Tabs.Tab>
          <Tabs.Tab
            value="budgets"
            leftSection={<DollarSign size={14} className="text-studio-muted" />}
          >
            Budgets
          </Tabs.Tab>
          <Tabs.Tab
            value="channels"
            leftSection={<Share2 size={14} className="text-studio-muted" />}
          >
            Channels & Connections
          </Tabs.Tab>
        </Tabs.List>

        {/* Tab 1: General */}
        <Tabs.Panel value="general">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
            <div>
              <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                <Building size={16} className="text-studio-accent" />
                Workspace General Information
              </Text>
              <Text size="xs" c="dimmed">
                Configure studio workspace branding, timezone, and primary production metadata.
              </Text>
            </div>

            <Divider color="dark.5" />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <TextInput
                label="Studio Workspace Name"
                defaultValue="Primas Digital Production Studio"
              />
              <Select
                label="Primary Timezone"
                defaultValue="UTC"
                data={[
                  { value: "UTC", label: "UTC (Coordinated Universal Time)" },
                  { value: "EST", label: "EST (Eastern Standard Time)" },
                  { value: "PST", label: "PST (Pacific Standard Time)" },
                ]}
              />
            </div>
          </Paper>
        </Tabs.Panel>

        {/* Tab 2: Production */}
        <Tabs.Panel value="production">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
            <div>
              <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                <Sliders size={16} className="text-studio-accent" />
                Production & Render Standards
              </Text>
              <Text size="xs" c="dimmed">
                Set baseline aspect ratios, framing rules, and continuity validation policy.
              </Text>
            </div>

            <Divider color="dark.5" />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Select
                label="Default Framing / Aspect Ratio"
                value={formData.defaultAspectRatio}
                onChange={(val) => val && handleChange("defaultAspectRatio", val as any)}
                data={[
                  { value: "16:9", label: "16:9 Landscape (Film/TV)" },
                  { value: "9:16", label: "9:16 Vertical (Mobile/Shorts)" },
                  { value: "1:1", label: "1:1 Square" },
                  { value: "4:3", label: "4:3 Academy Standard" },
                ]}
              />

              <Select
                label="Master Render Resolution"
                value={formData.defaultResolution}
                onChange={(val) => val && handleChange("defaultResolution", val as any)}
                data={[
                  { value: "1080p", label: "1080p Full HD" },
                  { value: "4K", label: "4K Cinema Ultra HD" },
                  { value: "720p", label: "720p Workprint Preview" },
                ]}
              />

              <Select
                label="Continuity Validation Policy"
                value={formData.continuityCheckSeverity}
                onChange={(val) => val && handleChange("continuityCheckSeverity", val as any)}
                data={[
                  { value: "strict", label: "Strict — Block generation on minor wardrobe/prop drift" },
                  { value: "moderate", label: "Moderate — Flag warnings, block on character shift" },
                  { value: "lenient", label: "Lenient — Log continuity notes without blocking" },
                ]}
              />
            </div>
          </Paper>
        </Tabs.Panel>

        {/* Tab 3: Models & Providers Keys & Connections */}
        <Tabs.Panel value="models">
          <div className="space-y-6">
            {/* Section 1: Default Production Models */}
            <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
              <div>
                <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                  <Cpu size={16} className="text-studio-accent" />
                  Production Models & Services
                </Text>
                <Text size="xs" c="dimmed">
                  Configure default models used for scriptwriting, world building, and media generation.
                </Text>
              </div>

              <Divider color="dark.5" />

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Select
                  label="Scriptwriting & Dialogue Model"
                  value={formData.scriptModel}
                  onChange={(val) => val && handleChange("scriptModel", val)}
                  data={[
                    { value: "gpt-4o", label: "OpenAI GPT-4o (Recommended)" },
                    { value: "claude-3-5-sonnet", label: "Anthropic Claude 3.5 Sonnet" },
                    { value: "gemini-1.5-pro", label: "Google Gemini 1.5 Pro" },
                  ]}
                />

                <Select
                  label="World Building Intelligence"
                  value={formData.worldBuildingModel}
                  onChange={(val) => val && handleChange("worldBuildingModel", val)}
                  data={[
                    { value: "claude-3-5-sonnet", label: "Anthropic Claude 3.5 Sonnet (Recommended)" },
                    { value: "gpt-4o", label: "OpenAI GPT-4o" },
                    { value: "gemini-1.5-pro", label: "Google Gemini 1.5 Pro" },
                  ]}
                />

                <Select
                  label="Character Voice Engine"
                  value={formData.characterVoiceEngine}
                  onChange={(val) => val && handleChange("characterVoiceEngine", val)}
                  data={[
                    { value: "elevenlabs-multilingual-v2", label: "ElevenLabs Multilingual v2" },
                    { value: "openai-tts-1-hd", label: "OpenAI TTS-1 HD" },
                  ]}
                />

                <Select
                  label="Video Generation Engine"
                  value={formData.videoGenModel}
                  onChange={(val) => val && handleChange("videoGenModel", val)}
                  data={[
                    { value: "runway-gen3-alpha", label: "Runway Gen-3 Alpha" },
                    { value: "luma-dream-machine", label: "Luma Dream Machine" },
                    { value: "pika-2", label: "Pika 2.0 Video Engine" },
                  ]}
                />
              </div>
            </Paper>

            {/* Section 2: Model Provider API Keys & Connection Testing */}
            <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
              <div>
                <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                  <Key size={16} className="text-studio-accent" />
                  Model Provider Keys & Connection Testing
                </Text>
                <Text size="xs" c="dimmed">
                  Enter provider keys and test endpoint connection status for each service.
                </Text>
              </div>

              {testResult && (
                <Alert
                  icon={testResult.success ? <CheckCircle2 size={16} /> : <AlertTriangle size={16} />}
                  title={`${testResult.provider} Connection Test`}
                  color={testResult.success ? "green" : "red"}
                  variant="light"
                  withCloseButton
                  onClose={() => setTestResult(null)}
                >
                  {testResult.message}
                </Alert>
              )}

              <Divider color="dark.5" />

              <div className="space-y-4">
                {[
                  {
                    key: "openai",
                    name: "OpenAI",
                    desc: "Scriptwriting, dialogue, and TTS synthesis",
                    field: "openaiKey" as const,
                  },
                  {
                    key: "anthropic",
                    name: "Anthropic",
                    desc: "World building, canon fact intelligence, and story logic",
                    field: "anthropicKey" as const,
                  },
                  {
                    key: "gemini",
                    name: "Google Gemini",
                    desc: "Long-context analysis and multimodal story processing",
                    field: "geminiKey" as const,
                  },
                  {
                    key: "elevenlabs",
                    name: "ElevenLabs",
                    desc: "Multilingual character voice synthesis",
                    field: "elevenLabsKey" as const,
                  },
                  {
                    key: "runway",
                    name: "Runway",
                    desc: "Gen-3 motion video generation",
                    field: "runwayKey" as const,
                  },
                  {
                    key: "midjourney",
                    name: "Midjourney / FLUX",
                    desc: "Character renders and concept storyboard panels",
                    field: "midjourneyKey" as const,
                  },
                  {
                    key: "replicate",
                    name: "Replicate",
                    desc: "Open-source AI video and image model hosting",
                    field: "replicateKey" as const,
                  },
                ].map((item) => {
                  const status = providerStatuses[item.key]?.status ?? "connected";
                  const isTesting = testingProvider === item.key;

                  return (
                    <Paper
                      key={item.key}
                      p="md"
                      radius="sm"
                      className="bg-studio-panel border border-studio-border space-y-3"
                    >
                      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                        <div>
                          <div className="flex items-center gap-2">
                            <Text fw={700} size="xs" c="white">
                              {item.name}
                            </Text>
                            <Badge
                              color={
                                status === "connected"
                                  ? "emerald"
                                  : status === "testing"
                                  ? "blue"
                                  : status === "invalid"
                                  ? "red"
                                  : "gray"
                              }
                              variant="light"
                              size="xs"
                            >
                              {status === "connected"
                                ? "● Connected"
                                : status === "testing"
                                ? "Testing..."
                                : status === "invalid"
                                ? "Connection Failed"
                                : "Disconnected"}
                            </Badge>
                          </div>
                          <Text size="xs" c="dimmed" mt={1}>
                            {item.desc}
                          </Text>
                        </div>

                        <Button
                          variant="outline"
                          color="terracotta"
                          size="xs"
                          loading={isTesting}
                          leftSection={!isTesting && <Zap size={14} />}
                          onClick={() =>
                            handleTestProviderConnection(item.key, item.name, formData[item.field] || "")
                          }
                        >
                          Test Connection
                        </Button>
                      </div>

                      <PasswordInput
                        size="xs"
                        placeholder={`Enter ${item.name} API Key`}
                        value={formData[item.field] || ""}
                        onChange={(e) => handleChange(item.field, e.currentTarget.value)}
                      />
                    </Paper>
                  );
                })}
              </div>
            </Paper>
          </div>
        </Tabs.Panel>

        {/* Tab 4: Storage */}
        <Tabs.Panel value="storage">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
            <div>
              <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                <HardDrive size={16} className="text-studio-accent" />
                Storage & Asset Retention
              </Text>
              <Text size="xs" c="dimmed">
                Configure cloud media storage providers and retention schedules for master video renders.
              </Text>
            </div>

            <Divider color="dark.5" />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Select
                label="Primary Cloud Asset Storage"
                defaultValue="s3"
                data={[
                  { value: "s3", label: "Amazon S3 Studio Vault" },
                  { value: "gcs", label: "Google Cloud Storage" },
                  { value: "cloudflare", label: "Cloudflare R2" },
                ]}
              />
              <Select
                label="Media Retention Policy"
                defaultValue="forever"
                data={[
                  { value: "forever", label: "Indefinite Archival (Recommended)" },
                  { value: "90days", label: "Purge Raw Workprints After 90 Days" },
                ]}
              />
            </div>
          </Paper>
        </Tabs.Panel>

        {/* Tab 5: Team */}
        <Tabs.Panel value="team">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
            <div>
              <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                <Users size={16} className="text-studio-accent" />
                Team Access & Role Governance
              </Text>
              <Text size="xs" c="dimmed">
                Manage producers, directors, operators, and distribution managers.
              </Text>
            </div>

            <Divider color="dark.5" />

            <Group justify="space-between" align="center">
              <div>
                <Text fw={600} size="xs" c="white">Studio Operators</Text>
                <Text size="xs" c="dimmed">3 members configured with executive rights</Text>
              </div>
              <Button size="xs" color="terracotta">Invite Team Member</Button>
            </Group>
          </Paper>
        </Tabs.Panel>

        {/* Tab 6: Budgets */}
        <Tabs.Panel value="budgets">
          <Paper p="xl" radius="md" className="bg-studio-card border border-studio-border space-y-6">
            <div>
              <Text fw={700} size="sm" c="white" className="flex items-center gap-2">
                <DollarSign size={16} className="text-studio-accent" />
                Budget Limits & Cost Guards
              </Text>
              <Text size="xs" c="dimmed">
                Establish spending caps per shot generation and maximum project budgets.
              </Text>
            </div>

            <Divider color="dark.5" />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <NumberInput
                label="Max Generation Cost Cap ($)"
                value={formData.maxCostPerGeneration}
                min={0.5}
                max={50}
                step={0.5}
                onChange={(val) => typeof val === "number" && handleChange("maxCostPerGeneration", val)}
              />
              <NumberInput
                label="Default Production Budget Cap ($)"
                value={formData.totalProjectBudgetCap}
                min={50}
                max={10000}
                step={50}
                onChange={(val) => typeof val === "number" && handleChange("totalProjectBudgetCap", val)}
              />
            </div>
          </Paper>
        </Tabs.Panel>

        {/* Tab 7: Channels & Connections */}
        <Tabs.Panel value="channels">
          <div className="space-y-6">
            {/* Header info */}
            <div>
              <Text fw={700} size="md" c="white" className="tracking-tight uppercase text-xs text-studio-accent mb-1">
                CHANNELS & CONNECTIONS
              </Text>
              <Text size="sm" c="dimmed">
                Connect the channels where your productions will be published.
              </Text>
            </div>

            {/* Connected Section */}
            <div className="space-y-3">
              <Text size="xs" fw={700} c="dimmed" className="tracking-wider uppercase">
                CONNECTED
              </Text>

              {connectedChannels.length > 0 ? (
                <div className="space-y-3">
                  {connectedChannels.map((channel) => (
                    <Paper
                      key={channel.id}
                      p="lg"
                      radius="md"
                      className="bg-studio-card border border-studio-border space-y-4"
                    >
                      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                        <div className="space-y-1">
                          <div className="flex items-center gap-3">
                            <Text fw={700} size="sm" c="white">
                              {providerLabels[channel.provider]?.title || channel.provider}
                            </Text>
                            <Badge
                              color="emerald"
                              variant="dot"
                              size="xs"
                              className="bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                            >
                              Connected
                            </Badge>
                          </div>
                          <Text size="xs" c="dimmed">
                            {channel.handle}
                          </Text>
                        </div>

                        <Button
                          variant="outline"
                          color="gray"
                          size="xs"
                          onClick={() => handleOpenManage(channel)}
                        >
                          Manage
                        </Button>
                      </div>

                      <Divider color="dark.5" />

                      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs">
                        <div>
                          <Text size="xs" c="dimmed">Publishing</Text>
                          <Text fw={600} c={channel.publishingEnabled ? "emerald" : "gray"}>
                            {channel.publishingEnabled ? "Enabled" : "Disabled"}
                          </Text>
                        </div>
                        <div>
                          <Text size="xs" c="dimmed">Auto-posting</Text>
                          <Text fw={600} c={channel.autoPostEnabled ? "emerald" : "gray"}>
                            {channel.autoPostEnabled ? "Enabled" : "Disabled"}
                          </Text>
                        </div>
                        <div>
                          <Text size="xs" c="dimmed">Default Format</Text>
                          <Text fw={600} c="white">
                            {channel.defaultFormat}
                          </Text>
                        </div>
                        <div>
                          <Text size="xs" c="dimmed">Default Publication</Text>
                          <Text fw={600} c="white">
                            {channel.defaultPublication}
                          </Text>
                        </div>
                      </div>
                    </Paper>
                  ))}
                </div>
              ) : (
                <Paper p="lg" radius="md" className="bg-studio-panel border border-studio-border text-center">
                  <Text size="xs" c="dimmed" fs="italic">
                    No distribution channels connected yet. Select an available channel below to begin setup.
                  </Text>
                </Paper>
              )}
            </div>

            {/* Available Channels Section */}
            <div className="space-y-3 pt-4">
              <Text size="xs" fw={700} c="dimmed" className="tracking-wider uppercase">
                AVAILABLE CHANNELS
              </Text>

              <Paper radius="md" className="bg-studio-card border border-studio-border divide-y divide-studio-border">
                {availableChannels.map((channel) => {
                  const meta = providerLabels[channel.provider];
                  return (
                    <div
                      key={channel.id}
                      className="p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:bg-studio-panel/50 transition-colors"
                    >
                      <div>
                        <Text fw={700} size="sm" c="white">
                          {meta?.title}
                        </Text>
                        <Text size="xs" c="dimmed" mt={1}>
                          {meta?.desc}
                        </Text>
                      </div>

                      <Button
                        color="terracotta"
                        size="xs"
                        onClick={() => handleStartOAuthConnect(channel.provider)}
                      >
                        Connect
                      </Button>
                    </div>
                  );
                })}

                {availableChannels.length === 0 && (
                  <div className="p-5 text-center">
                    <Text size="xs" c="emerald" fw={600}>
                      All available distribution channels are currently connected.
                    </Text>
                  </div>
                )}
              </Paper>
            </div>
          </div>
        </Tabs.Panel>
      </Tabs>

      {/* MANAGE CHANNEL MODAL */}
      <Modal
        opened={isManageModalOpen}
        onClose={() => setIsManageModalOpen(false)}
        title={
          <div className="flex items-center gap-2">
            <Text fw={700} size="sm" c="white">
              Manage {managingChannel ? providerLabels[managingChannel.provider]?.title : ""} Connection
            </Text>
            <Badge color="emerald" variant="dot" size="xs">
              Connected
            </Badge>
          </div>
        }
        centered
        size="lg"
        classNames={{
          content: "bg-studio-card border border-studio-border text-white",
          header: "bg-studio-card border-b border-studio-border text-white",
        }}
      >
        {managingChannel && (
          <div className="space-y-6 pt-2">
            {/* Account Info */}
            <div className="bg-studio-panel p-4 rounded border border-studio-border flex justify-between items-center">
              <div>
                <Text fw={700} size="sm" c="white">
                  {managingChannel.handle}
                </Text>
                <Text size="xs" c="dimmed">
                  Account: {managingChannel.accountName}
                </Text>
              </div>
              <Badge color="emerald" variant="light" size="xs">
                ● Connected
              </Badge>
            </div>

            {/* Permissions */}
            <div className="space-y-2">
              <Text fw={600} size="xs" c="dimmed" className="uppercase tracking-wider">
                Permissions Granted
              </Text>
              <div className="bg-studio-panel/50 p-3 rounded border border-studio-border space-y-1.5 text-xs">
                {managingChannel.permissions.map((perm) => (
                  <div key={perm} className="flex items-center gap-2 text-emerald-400">
                    <Check size={14} />
                    <span>{perm}</span>
                  </div>
                ))}
              </div>
            </div>

            <Divider color="dark.5" />

            {/* Publishing Controls */}
            <div className="space-y-4">
              <Text fw={600} size="xs" c="white">
                Publishing Settings
              </Text>

              <Group justify="space-between">
                <div>
                  <Text size="xs" fw={600} c="white">
                    Allow automatic publishing
                  </Text>
                  <Text size="xs" c="dimmed">
                    Enable auto-posting for approved episodes in publishing plans
                  </Text>
                </div>
                <Switch
                  color="terracotta"
                  checked={managingChannel.autoPostEnabled}
                  onChange={(e) =>
                    setManagingChannel({
                      ...managingChannel,
                      autoPostEnabled: e.currentTarget.checked,
                    })
                  }
                />
              </Group>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <Select
                  label="Default Format"
                  value={managingChannel.defaultFormat}
                  onChange={(val) =>
                    val &&
                    setManagingChannel({
                      ...managingChannel,
                      defaultFormat: val as any,
                    })
                  }
                  data={["Vertical 9:16", "Horizontal 16:9", "Square 1:1"]}
                />

                <Select
                  label="Default Caption"
                  value={managingChannel.defaultCaption}
                  onChange={(val) =>
                    val &&
                    setManagingChannel({
                      ...managingChannel,
                      defaultCaption: val as any,
                    })
                  }
                  data={["Use production caption", "Custom template"]}
                />

                <Select
                  label="Default Publication"
                  value={managingChannel.defaultPublication}
                  onChange={(val) =>
                    val &&
                    setManagingChannel({
                      ...managingChannel,
                      defaultPublication: val as any,
                    })
                  }
                  data={["Immediately", "Draft", "Scheduled Window"]}
                />
              </div>
            </div>

            <Divider color="dark.5" />

            {/* Save & Danger Zone */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-2">
              <Button
                color="red"
                variant="subtle"
                size="xs"
                leftSection={<Trash2 size={14} />}
                onClick={handleDisconnectChannel}
              >
                Disconnect Account
              </Button>

              <Group gap="xs">
                <Button
                  variant="outline"
                  color="gray"
                  size="xs"
                  onClick={() => setIsManageModalOpen(false)}
                >
                  Cancel
                </Button>
                <Button color="terracotta" size="xs" onClick={handleSaveChannelManagement}>
                  Save Changes
                </Button>
              </Group>
            </div>
          </div>
        )}
      </Modal>

      {/* OAUTH FLOW SIMULATION MODAL */}
      <Modal
        opened={isOAuthModalOpen}
        onClose={() => setIsOAuthModalOpen(false)}
        title={
          <Text fw={700} size="sm" c="white">
            Connect Channel
          </Text>
        }
        centered
        size="md"
        classNames={{
          content: "bg-studio-card border border-studio-border text-white",
          header: "bg-studio-card border-b border-studio-border text-white",
        }}
      >
        {connectingProvider && (
          <div className="space-y-5 pt-2">
            <Alert icon={<Shield size={16} />} color="blue" variant="outline">
              <Text size="xs" c="dimmed">
                Sign in with your <strong>{providerLabels[connectingProvider].title}</strong> account to grant publishing access.
              </Text>
            </Alert>

            <div className="bg-studio-panel p-4 rounded border border-studio-border text-xs text-studio-muted">
              <Text fw={600} c="white" mb={1}>Account Security:</Text>
              Your connection is secure. We only request permissions necessary to publish approved content and check publication status.
            </div>

            <div className="flex justify-end gap-3 pt-3 border-t border-studio-border">
              <Button variant="outline" color="gray" size="xs" onClick={() => setIsOAuthModalOpen(false)}>
                Cancel
              </Button>
              <Button
                color="terracotta"
                size="xs"
                leftSection={<ExternalLink size={14} />}
                onClick={handleCompleteOAuthConnect}
              >
                Authorize & Connect Channel
              </Button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
