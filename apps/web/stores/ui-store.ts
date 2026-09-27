import { create } from 'zustand'

interface UIState {
  isSidebarExpanded: boolean
  isCommandPaletteOpen: boolean
  toggleSidebar: () => void
  setCommandPaletteOpen: (open: boolean) => void
}

export const useUIStore = create<UIState>((set) => ({
  isSidebarExpanded: true,
  isCommandPaletteOpen: false,
  toggleSidebar: () => set((state) => ({ isSidebarExpanded: !state.isSidebarExpanded })),
  setCommandPaletteOpen: (open) => set({ isCommandPaletteOpen: open }),
}))
