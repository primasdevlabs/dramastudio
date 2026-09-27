import { create } from 'zustand'

export type ProductionMode = 'monitored' | 'autonomous'

interface StudioState {
  activeProjectId: string
  productionMode: ProductionMode
  activeEpisodeId: string | null
  setActiveProject: (id: string) => void
  setProductionMode: (mode: ProductionMode) => void
  setActiveEpisode: (id: string | null) => void
}

export const useStudioStore = create<StudioState>((set) => ({
  activeProjectId: 'prj-01',
  productionMode: 'monitored',
  activeEpisodeId: 'ep-12',
  setActiveProject: (id) => set({ activeProjectId: id }),
  setProductionMode: (mode) => set({ productionMode: mode }),
  setActiveEpisode: (id) => set({ activeEpisodeId: id }),
}))
