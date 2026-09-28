import { create } from "zustand";
import { persist } from "zustand/middleware";
import { GlobalSettings, ProjectSettingsOverride } from "@/lib/api/settings-types";

export const DEFAULT_GLOBAL_SETTINGS: GlobalSettings = {
  // AI Config
  scriptModel: "gpt-4o",
  worldBuildingModel: "claude-3-5-sonnet",
  characterVoiceEngine: "elevenlabs-multilingual-v2",
  imageGenModel: "midjourney-v6.1",
  videoGenModel: "runway-gen3-alpha",
  creativeTemperature: 0.7,

  // Production Defaults
  defaultAspectRatio: "16:9",
  defaultResolution: "1080p",
  maxCostPerGeneration: 2.5,
  totalProjectBudgetCap: 500,
  currency: "USD",
  autoApproveLowCostShots: false,
  continuityCheckSeverity: "strict",

  // API Keys & Model Provider Connections
  openaiKey: "sk-proj-••••••••••••••••••••",
  anthropicKey: "sk-ant-••••••••••••••••••••",
  geminiKey: "AIzaSy-••••••••••••••••••••",
  elevenLabsKey: "el-live-••••••••••••••••••••",
  runwayKey: "rw-secret-••••••••••••••••••••",
  replicateKey: "r8-token-••••••••••••••••••••",
  midjourneyKey: "mj-api-••••••••••••••••••••",
  providerStatuses: {
    openai: { status: "connected", lastTested: "Just now" },
    anthropic: { status: "connected", lastTested: "Just now" },
    gemini: { status: "connected", lastTested: "Just now" },
    elevenlabs: { status: "connected", lastTested: "Just now" },
    runway: { status: "connected", lastTested: "Just now" },
    replicate: { status: "connected", lastTested: "Just now" },
    midjourney: { status: "connected", lastTested: "Just now" },
  },

  // Channels & Connections (OAuth Distribution System)
  facebook: {
    id: "facebook",
    provider: "facebook",
    accountName: "Primas Studio Page",
    handle: "Primas Studio Page",
    connected: true,
    publishingEnabled: true,
    autoPostEnabled: true,
    defaultFormat: "Horizontal 16:9",
    defaultCaption: "Use production caption",
    defaultPublication: "Immediately",
    permissions: [
      "Publish content",
      "Read publishing status",
      "Read account information",
    ],
    capabilities: {
      formats: ["video", "image"],
      supportsScheduling: true,
      supportsThumbnails: true,
      maxVideoDurationSeconds: 14400,
    },
  },
  instagram: {
    id: "instagram",
    provider: "instagram",
    accountName: "Primas Studio",
    handle: "@studio",
    connected: true,
    publishingEnabled: true,
    autoPostEnabled: true,
    defaultFormat: "Vertical 9:16",
    defaultCaption: "Use production caption",
    defaultPublication: "Immediately",
    permissions: [
      "Publish content",
      "Read publishing status",
      "Read account information",
    ],
    capabilities: {
      formats: ["reel", "image", "carousel"],
      supportsScheduling: true,
      supportsThumbnails: true,
      maxVideoDurationSeconds: 90,
    },
  },
  tiktok: {
    id: "tiktok",
    provider: "tiktok",
    accountName: "DramaStudio Creator",
    handle: "@dramastudio_official",
    connected: false,
    publishingEnabled: false,
    autoPostEnabled: false,
    defaultFormat: "Vertical 9:16",
    defaultCaption: "Use production caption",
    defaultPublication: "Immediately",
    permissions: [
      "Publish content",
      "Read publishing status",
      "Read account information",
    ],
    capabilities: {
      formats: ["video"],
      supportsScheduling: true,
      supportsThumbnails: false,
      maxVideoDurationSeconds: 600,
    },
  },
  youtube: {
    id: "youtube",
    provider: "youtube",
    accountName: "DramaStudio Official",
    handle: "DramaStudio Shorts",
    connected: false,
    publishingEnabled: false,
    autoPostEnabled: false,
    defaultFormat: "Horizontal 16:9",
    defaultCaption: "Use production caption",
    defaultPublication: "Immediately",
    permissions: [
      "Publish content",
      "Read publishing status",
      "Read account information",
    ],
    capabilities: {
      formats: ["video", "shorts"],
      supportsScheduling: true,
      supportsThumbnails: true,
      maxVideoDurationSeconds: 43200,
    },
  },
  autoPublishOnApproval: false,
  preferredPostingWindow: "18:00 - 21:00 EST",
};

interface SettingsStoreState {
  globalSettings: GlobalSettings;
  projectOverrides: Record<string, ProjectSettingsOverride>;

  updateGlobalSettings: (settings: Partial<GlobalSettings>) => void;
  resetGlobalSettings: () => void;
  
  toggleProjectInheritance: (projectId: string, inherit: boolean) => void;
  updateProjectOverride: (projectId: string, overrides: Partial<GlobalSettings>) => void;
  updateProjectMetadata: (projectId: string, metadata: { name?: string; logline?: string; genre?: string; targetAudience?: string }) => void;
  
  getEffectiveSettings: (projectId?: string) => GlobalSettings;
  getProjectOverrideState: (projectId: string) => ProjectSettingsOverride;
}

export const useSettingsStore = create<SettingsStoreState>()(
  persist(
    (set, get) => ({
      globalSettings: DEFAULT_GLOBAL_SETTINGS,
      projectOverrides: {},

      updateGlobalSettings: (newSettings) =>
        set((state) => ({
          globalSettings: { ...state.globalSettings, ...newSettings },
        })),

      resetGlobalSettings: () =>
        set({ globalSettings: DEFAULT_GLOBAL_SETTINGS }),

      toggleProjectInheritance: (projectId, inherit) =>
        set((state) => {
          const current = state.projectOverrides[projectId] || {
            inheritGlobalDefaults: true,
            overrides: {},
          };
          return {
            projectOverrides: {
              ...state.projectOverrides,
              [projectId]: { ...current, inheritGlobalDefaults: inherit },
            },
          };
        }),

      updateProjectOverride: (projectId, newOverrides) =>
        set((state) => {
          const current = state.projectOverrides[projectId] || {
            inheritGlobalDefaults: true,
            overrides: {},
          };
          return {
            projectOverrides: {
              ...state.projectOverrides,
              [projectId]: {
                ...current,
                overrides: { ...current.overrides, ...newOverrides },
              },
            },
          };
        }),

      updateProjectMetadata: (projectId, metadata) =>
        set((state) => {
          const current = state.projectOverrides[projectId] || {
            inheritGlobalDefaults: true,
            overrides: {},
          };
          return {
            projectOverrides: {
              ...state.projectOverrides,
              [projectId]: {
                ...current,
                ...metadata,
              },
            },
          };
        }),

      getProjectOverrideState: (projectId) => {
        const state = get();
        return (
          state.projectOverrides[projectId] || {
            inheritGlobalDefaults: true,
            overrides: {},
          }
        );
      },

      getEffectiveSettings: (projectId) => {
        const state = get();
        if (!projectId) return state.globalSettings;
        const projectState = state.projectOverrides[projectId];
        if (!projectState || projectState.inheritGlobalDefaults) {
          return state.globalSettings;
        }
        return {
          ...state.globalSettings,
          ...projectState.overrides,
        };
      },
    }),
    {
      name: "dramastudio-settings-storage",
    }
  )
);
