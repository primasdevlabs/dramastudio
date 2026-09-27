import Link from 'next/link'
import { Film, Bot, Sparkles, Layers, Cpu, ArrowRight } from 'lucide-react'

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-studio-bg flex flex-col justify-between p-8">
      <header className="max-w-7xl mx-auto w-full flex justify-between items-center py-6 border-b border-studio-border">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-studio-accent to-studio-violet flex items-center justify-center glow-cyan">
            <Film className="w-5 h-5 text-white" />
          </div>
          <span className="font-bold text-xl tracking-tight text-white">DramaStudio</span>
        </div>
        <div className="flex items-center gap-6">
          <Link href="/pricing" className="text-sm text-studio-muted hover:text-white transition">Pricing</Link>
          <Link href="/about" className="text-sm text-studio-muted hover:text-white transition">About</Link>
          <Link href="/projects" className="px-5 py-2.5 rounded-lg bg-studio-accent text-slate-950 font-semibold text-sm hover:bg-cyan-400 transition glow-cyan flex items-center gap-2">
            Enter Studio <ArrowRight className="w-4 h-4" />
          </Link>
        </div>
      </header>

      <main className="max-w-5xl mx-auto text-center py-20">
        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass-panel text-studio-accent text-xs font-semibold uppercase tracking-wider mb-8">
          <Sparkles className="w-3.5 h-3.5" /> AI-Native Serialized Drama Production
        </div>
        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-6 bg-gradient-to-r from-white via-slate-200 to-studio-muted bg-clip-text text-transparent">
          AI Creates Content.<br/>The Studio Owns Reality.
        </h1>
        <p className="text-lg md:text-xl text-studio-muted max-w-2xl mx-auto mb-10 leading-relaxed">
          Orchestrate specialized AI model providers through an autonomous Lead Director agent. Turn story ideas into 150-episode serialized dramas with persistent story continuity.
        </p>

        <div className="flex flex-wrap justify-center gap-4 mb-16">
          <Link href="/projects" className="px-8 py-4 rounded-xl bg-gradient-to-r from-studio-accent to-studio-violet text-white font-bold text-base hover:opacity-90 transition glow-violet flex items-center gap-3">
            Open Studio Workspace <ArrowRight className="w-5 h-5" />
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-left max-w-4xl mx-auto">
          <div className="p-6 rounded-2xl glass-panel">
            <Bot className="w-8 h-8 text-studio-accent mb-4" />
            <h3 className="font-bold text-white mb-2">Lead Director Agent</h3>
            <p className="text-sm text-studio-muted leading-normal">Autonomous orchestration of story, visual, audio, editing, and QA sub-agents.</p>
          </div>
          <div className="p-6 rounded-2xl glass-panel">
            <Cpu className="w-8 h-8 text-studio-violet mb-4" />
            <h3 className="font-bold text-white mb-2">Capability Model Policy</h3>
            <p className="text-sm text-studio-muted leading-normal">Provider-agnostic model registry for script, dialogue, storyboard, Wan video, and TTS.</p>
          </div>
          <div className="p-6 rounded-2xl glass-panel">
            <Layers className="w-8 h-8 text-studio-pink mb-4" />
            <h3 className="font-bold text-white mb-2">Persistent Story Canon</h3>
            <p className="text-sm text-studio-muted leading-normal">Fact knowledge graph & automated multi-faceted continuity validation across seasons.</p>
          </div>
        </div>
      </main>

      <footer className="max-w-7xl mx-auto w-full text-center text-xs text-studio-muted border-t border-studio-border pt-6">
        Copyright © 2026 DramaStudio. All rights reserved.
      </footer>
    </div>
  )
}
