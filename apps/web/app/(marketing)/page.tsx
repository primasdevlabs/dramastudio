import Link from 'next/link'
import { Film, Bot, Layers, Cpu, ArrowRight, Clapperboard } from 'lucide-react'

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-studio-bg flex flex-col justify-between p-8">
      <header className="max-w-7xl mx-auto w-full flex justify-between items-center py-6 border-b border-studio-border">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-md bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
            <Film className="w-5 h-5" />
          </div>
          <span className="font-bold text-xl tracking-tight text-white">DramaStudio</span>
        </div>
        <div className="flex items-center gap-6">
          <Link href="/pricing" className="text-sm text-studio-muted hover:text-white transition">Pricing</Link>
          <Link href="/about" className="text-sm text-studio-muted hover:text-white transition">About</Link>
          <Link href="/projects" className="px-5 py-2.5 rounded bg-studio-accent text-white font-medium text-sm hover:bg-studio-accent-dark transition flex items-center gap-2">
            Enter Studio <ArrowRight className="w-4 h-4" />
          </Link>
        </div>
      </header>

      <main className="max-w-5xl mx-auto text-center py-20">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded bg-studio-panel border border-studio-border text-studio-accent text-xs font-mono uppercase tracking-wider mb-8">
          <Clapperboard className="w-3.5 h-3.5" /> Serialized Drama Production Workstation
        </div>
        <h1 className="text-4xl md:text-6xl font-bold tracking-tight text-white mb-6">
          Automated Production.<br/>Persistent Story Canon.
        </h1>
        <p className="text-lg md:text-xl text-studio-muted max-w-2xl mx-auto mb-10 leading-relaxed">
          Orchestrate specialized model backends through a Lead Director console. Turn story premises into 150-episode serialized dramas with persistent canon continuity.
        </p>

        <div className="flex flex-wrap justify-center gap-4 mb-16">
          <Link href="/projects" className="px-8 py-4 rounded bg-studio-accent text-white font-bold text-base hover:bg-studio-accent-dark transition flex items-center gap-3">
            Open Studio Console <ArrowRight className="w-5 h-5" />
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-left max-w-4xl mx-auto">
          <div className="p-6 rounded-md bg-studio-card border border-studio-border">
            <Bot className="w-8 h-8 text-studio-accent mb-4" />
            <h3 className="font-bold text-white mb-2">Lead Director Console</h3>
            <p className="text-sm text-studio-muted leading-normal">Orchestration of story, visual, audio, assembly, and continuity roles.</p>
          </div>
          <div className="p-6 rounded-md bg-studio-card border border-studio-border">
            <Cpu className="w-8 h-8 text-amber-400 mb-4" />
            <h3 className="font-bold text-white mb-2">Capability Model Policy</h3>
            <p className="text-sm text-studio-muted leading-normal">Provider-agnostic model registry for script, dialogue, storyboard, video shot, and voice track.</p>
          </div>
          <div className="p-6 rounded-md bg-studio-card border border-studio-border">
            <Layers className="w-8 h-8 text-emerald-400 mb-4" />
            <h3 className="font-bold text-white mb-2">Persistent Story Canon</h3>
            <p className="text-sm text-studio-muted leading-normal">Canon fact knowledge graph & automated multi-faceted continuity validation across seasons.</p>
          </div>
        </div>
      </main>

      <footer className="max-w-7xl mx-auto w-full text-center text-xs text-studio-muted border-t border-studio-border pt-6">
        Copyright © 2026 DramaStudio. All rights reserved.
      </footer>
    </div>
  )
}
