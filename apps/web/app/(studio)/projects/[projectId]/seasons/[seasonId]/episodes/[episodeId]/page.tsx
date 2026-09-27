'use client'

import { useState } from 'react'
import { Sparkles, FileText, Image, Video, Mic, Film, ShieldAlert, Bot, Play, CheckCircle2 } from 'lucide-react'
import { useWorkspaceStore } from '@/stores/workspace-store'

export default function EpisodeProductionWorkspace() {
  const { activeTab, setActiveTab } = useWorkspaceStore()

  const capabilities = [
    { key: 'script', label: 'Script', icon: FileText },
    { key: 'storyboard', label: 'Storyboard', icon: Image },
    { key: 'shots', label: 'Wan Video Shots', icon: Video },
    { key: 'audio', label: 'Voice & SFX', icon: Mic },
    { key: 'assembly', label: 'FFmpeg Assembly', icon: Film },
    { key: 'continuity', label: 'Continuity & QA', icon: ShieldAlert },
  ]

  return (
    <div className="space-y-6">
      {/* Workspace Header Bar */}
      <div className="flex justify-between items-center bg-studio-card p-4 rounded-2xl border border-studio-border">
        <div>
          <div className="flex items-center gap-2 text-xs text-studio-muted">
            <span>Season 1</span> • <span>Episode 12</span>
          </div>
          <h1 className="text-xl font-bold text-white">Scene 04: The Neon Alley Confrontation</h1>
        </div>

        <div className="flex items-center gap-3">
          <button className="px-4 py-2 bg-studio-panel border border-studio-border text-xs font-semibold text-white rounded-lg flex items-center gap-2 hover:bg-studio-border">
            <Bot className="w-4 h-4 text-studio-accent" /> Trigger Lead Director Run
          </button>
          <button className="px-5 py-2 bg-gradient-to-r from-studio-accent to-studio-violet text-white font-bold text-xs rounded-lg flex items-center gap-2 glow-cyan">
            <Play className="w-4 h-4 fill-current" /> Render Episode Assembly
          </button>
        </div>
      </div>

      {/* Workspace Sub-Navigation Tabs */}
      <div className="flex border-b border-studio-border gap-2">
        {capabilities.map((item) => {
          const Icon = item.icon
          const isActive = activeTab === item.key
          return (
            <button
              key={item.key}
              onClick={() => setActiveTab(item.key as any)}
              className={`px-4 py-2.5 text-xs font-semibold flex items-center gap-2 border-b-2 transition ${
                isActive
                  ? 'border-studio-accent text-studio-accent bg-studio-card/50'
                  : 'border-transparent text-studio-muted hover:text-white'
              }`}
            >
              <Icon className="w-4 h-4" /> {item.label}
            </button>
          )}
        )}
      </div>

      {/* Main Workspace Workspace Panels */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left 2 Cols: Main Artifact Workspace */}
        <div className="lg:col-span-2 space-y-4">
          <div className="p-6 rounded-2xl glass-panel space-y-4">
            <div className="flex justify-between items-center border-b border-studio-border pb-3">
              <h3 className="font-bold text-white text-sm">Active Artifact: Storyboard Panel 07</h3>
              <span className="px-2.5 py-1 rounded bg-studio-accent/10 text-studio-accent text-[10px] font-bold">
                Capability: storyboard (Model: Vision-XL)
              </span>
            </div>

            {/* Visual Storyboard Frame */}
            <div className="aspect-video rounded-xl bg-studio-bg border border-studio-border flex flex-col items-center justify-center text-center p-6 relative overflow-hidden group">
              <div className="w-16 h-16 rounded-full bg-studio-panel flex items-center justify-center mb-3 text-studio-accent">
                <Video className="w-8 h-8" />
              </div>
              <p className="text-xs font-semibold text-white mb-1">Shot 07: Medium Close-up on Sarah Croft</p>
              <p className="text-[11px] text-studio-muted max-w-md">"Sarah reaches for her encrypted datapad as neon shadows reflect across her trench coat."</p>
              <div className="absolute bottom-3 right-3 px-2 py-1 rounded bg-black/60 backdrop-blur text-[10px] text-emerald-400 font-semibold">
                Duration: 4.2s • Wan 2.1 Video Ready
              </div>
            </div>

            {/* Prompt & Context Controls */}
            <div className="space-y-2">
              <label className="block text-[11px] font-semibold text-studio-muted">Specialized Model Context Payload (Resolved via Model Policy)</label>
              <div className="p-3 rounded-lg bg-studio-bg border border-studio-border text-[11px] text-studio-muted space-y-1">
                <div><span className="text-studio-accent">character:</span> Sarah Croft (Wardrobe: Cyberpunk Trench Coat)</div>
                <div><span className="text-studio-violet">location:</span> Neo-London Alley 14B</div>
                <div><span className="text-emerald-400">continuity_constraint:</span> Sarah holds injury on left arm</div>
              </div>
            </div>
          </div>
        </div>

        {/* Right Col: Live Director & Continuity Feed */}
        <div className="space-y-4">
          {/* Realtime Lead Director Activity */}
          <div className="p-4 rounded-2xl glass-panel space-y-3">
            <div className="flex items-center gap-2 border-b border-studio-border pb-2">
              <Bot className="w-4 h-4 text-studio-accent" />
              <h4 className="text-xs font-bold text-white">Lead Director Live Activity</h4>
            </div>
            <div className="space-y-2 text-[11px]">
              <div className="p-2.5 rounded-lg bg-studio-panel border border-studio-border flex gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
                <div>
                  <div className="font-semibold text-white">Script & Dialogue Approved</div>
                  <div className="text-[10px] text-studio-muted">Dialogue Model B generated Scene 4 script</div>
                </div>
              </div>
              <div className="p-2.5 rounded-lg bg-studio-panel border border-studio-border flex gap-2">
                <Sparkles className="w-4 h-4 text-studio-accent shrink-0 mt-0.5" />
                <div>
                  <div className="font-semibold text-white">Generating Storyboard Shots</div>
                  <div className="text-[10px] text-studio-muted">Resolved capability &apos;storyboard&apos; &rarr; Vision Provider X</div>
                </div>
              </div>
            </div>
          </div>

          {/* Continuity Issue Alert Box */}
          <div className="p-4 rounded-2xl bg-amber-500/10 border border-amber-500/20 space-y-2">
            <div className="flex items-center gap-2 text-amber-400 text-xs font-bold">
              <ShieldAlert className="w-4 h-4" /> 1 Continuity Issue Flagged
            </div>
            <p className="text-[11px] text-amber-200/80 leading-normal">
              Sarah's wardrobe in Scene 4 does not match Episode 11 concluding state. Continuity Agent recommends locking trench coat variant.
            </p>
            <button className="w-full py-1.5 rounded bg-amber-500 text-slate-950 text-[11px] font-bold">
              Auto-Resolve Continuity
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
