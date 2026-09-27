import { create } from "zustand";

interface StudioState {
  activeProjectId: string | null;
  activeEpisodeId: string | null;
  activeTab: string;
  isSidebarOpen: boolean;
  isCommandPaletteOpen: boolean;
  setActiveProjectId: (id: string | null) => void;
  setActiveEpisodeId: (id: string | null) => void;
  setActiveTab: (tab: string) => void;
  toggleSidebar: () => void;
  toggleCommandPalette: () => void;
  setCommandPaletteOpen: (open: boolean) => void;
}

export const useStudioStore = create<StudioState>((set) => ({
  activeProjectId: null,
  activeEpisodeId: null,
  activeTab: "overview",
  isSidebarOpen: true,
  isCommandPaletteOpen: false,
  setActiveProjectId: (id) => set({ activeProjectId: id }),
  setActiveEpisodeId: (id) => set({ activeEpisodeId: id }),
  setActiveTab: (tab) => set({ activeTab: tab }),
  toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),
  toggleCommandPalette: () =>
    set((state) => ({ isCommandPaletteOpen: !state.isCommandPaletteOpen })),
  setCommandPaletteOpen: (open) => set({ isCommandPaletteOpen: open }),
}));
