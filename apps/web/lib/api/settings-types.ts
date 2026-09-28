export interface StudioAIConfig {
  scriptModel: string;
  worldBuildingModel: string;
  characterVoiceEngine: string;
  imageGenModel: string;
  videoGenModel: string;
  creativeTemperature: number;
}

export interface StudioProductionDefaults {
  defaultAspectRatio: "16:9" | "9:16" | "1:1" | "4:3";
  defaultResolution: "720p" | "1080p" | "4K";
  maxCostPerGeneration: number;
  totalProjectBudgetCap: number;
  currency: string;
  autoApproveLowCostShots: boolean;
  continuityCheckSeverity: "strict" | "moderate" | "lenient";
}

export interface ProviderConnectionStatus {
  status: "connected" | "disconnected" | "testing" | "invalid";
  lastTested?: string;
  errorMessage?: string;
}

export interface StudioApiKeys {
  openaiKey: string;
  anthropicKey: string;
  geminiKey: string;
  elevenLabsKey: string;
  runwayKey: string;
  replicateKey: string;
  midjourneyKey: string;
  providerStatuses?: Record<string, ProviderConnectionStatus>;
}

export interface ChannelCapabilities {
  formats: ("video" | "reel" | "shorts" | "image" | "carousel")[];
  supportsScheduling: boolean;
  supportsThumbnails: boolean;
  maxVideoDurationSeconds: number;
}

export interface ChannelConnectionDetail {
  id: string;
  provider: "facebook" | "instagram" | "tiktok" | "youtube";
  accountName: string;
  handle: string;
  connected: boolean;
  publishingEnabled: boolean;
  autoPostEnabled: boolean;
  defaultFormat: "Vertical 9:16" | "Horizontal 16:9" | "Square 1:1";
  defaultCaption: "Use production caption" | "Custom template";
  defaultPublication: "Immediately" | "Draft" | "Scheduled Window";
  permissions: string[];
  capabilities: ChannelCapabilities;
}

export interface StudioSocialChannels {
  facebook: ChannelConnectionDetail;
  instagram: ChannelConnectionDetail;
  tiktok: ChannelConnectionDetail;
  youtube: ChannelConnectionDetail;
  autoPublishOnApproval: boolean;
  preferredPostingWindow: string;
}

export interface PublishingPlanChannelItem {
  channelId: string;
  channelName: string;
  enabled: boolean;
  format: string;
  publicationTime: string;
  status: "scheduled" | "publishing" | "published" | "failed" | "draft";
  failureReason?: string;
}

export interface PublishingPlan {
  id: string;
  episodeId: string;
  episodeTitle: string;
  scheduledTime: string;
  channels: PublishingPlanChannelItem[];
}

export type GlobalSettings = StudioAIConfig &
  StudioProductionDefaults &
  StudioApiKeys &
  StudioSocialChannels;

export interface ProjectSettingsOverride {
  inheritGlobalDefaults: boolean;
  overrides: Partial<GlobalSettings>;
  name?: string;
  logline?: string;
  genre?: string;
  targetAudience?: string;
}

