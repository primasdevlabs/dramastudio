import { create } from 'zustand'

interface WorkspaceState {
  activeTab: 'script' | 'scenes' | 'storyboard' | 'shots' | 'assets' | 'audio' | 'assembly' | 'continuity' | 'qa'
  selectedSceneId: string | null
  selectedShotId: string | null
  setActiveTab: (tab: WorkspaceState['activeTab']) => void
  setSelectedScene: (id: string | null) => void
  setSelectedShot: (id: string | null) => void
}

export const useWorkspaceStore = create<WorkspaceState>((set) => ({
  activeTab: 'storyboard',
  selectedSceneId: 'scene-04',
  selectedShotId: 'shot-07',
  setActiveTab: (tab) => set({ activeTab: tab }),
  setSelectedScene: (id) => set({ selectedSceneId: id }),
  setSelectedShot: (id) => set({ selectedShotId: id }),
}))
