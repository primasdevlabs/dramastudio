import Link from 'next/link'
import { Clapperboard, Sparkles, Users, ShieldCheck, Play, ArrowRight } from 'lucide-react'

export default function ProjectDashboardPage({ params }: { params: { projectId: string } }) {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-white">Shadows of Neo-London</h1>
          <p className="text-xs text-studio-muted">Series Overview & Autonomous Production Pipeline</p>
        </div>
        <Link href={`/projects/${params.projectId}/seasons/season-01/episodes/ep-12`} className="px-5 py-2.5 bg-gradient-to-r from-studio-accent to-studio-violet text-white font-bold text-xs rounded-xl shadow glow-cyan flex items-center gap-2">
          <Play className="w-4 h-4 fill-current" /> Open Episode Production Workspace
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl glass-panel">
          <div className="text-[10px] font-semibold text-studio-muted uppercase">Series Bible State</div>
          <div className="text-xl font-bold text-white mt-1">Canonical (v1.4)</div>
        </div>
        <div className="p-4 rounded-xl glass-panel">
          <div className="text-[10px] font-semibold text-studio-muted uppercase">Active Characters</div>
          <div className="text-xl font-bold text-studio-accent mt-1">8 Core Profiles</div>
        </div>
        <div className="p-4 rounded-xl glass-panel">
          <div className="text-[10px] font-semibold text-studio-muted uppercase">Canon Facts</div>
          <div className="text-xl font-bold text-studio-violet mt-1">193 Established</div>
        </div>
        <div className="p-4 rounded-xl glass-panel">
          <div className="text-[10px] font-semibold text-studio-muted uppercase">Continuity Status</div>
          <div className="text-xl font-bold text-emerald-400 mt-1">98.4% Clean</div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link href={`/projects/${params.projectId}/story`} className="p-6 rounded-2xl glass-panel hover:border-studio-accent transition">
          <Sparkles className="w-6 h-6 text-studio-accent mb-3" />
          <h3 className="font-bold text-white mb-1">Story & Canon Engine</h3>
          <p className="text-xs text-studio-muted">Manage Series Bible, Story Graph, and Fact Knowledge Base.</p>
        </Link>

        <Link href={`/projects/${params.projectId}/seasons/season-01/episodes/ep-12`} className="p-6 rounded-2xl glass-panel hover:border-studio-violet transition">
          <Clapperboard className="w-6 h-6 text-studio-violet mb-3" />
          <h3 className="font-bold text-white mb-1">Episode Production Workspace</h3>
          <p className="text-xs text-studio-muted">Script, Storyboard, Wan Video, Voice & Assembly pipeline.</p>
        </Link>

        <Link href={`/projects/${params.projectId}/continuity`} className="p-6 rounded-2xl glass-panel hover:border-emerald-400 transition">
          <ShieldCheck className="w-6 h-6 text-emerald-400 mb-3" />
          <h3 className="font-bold text-white mb-1">Continuity Engine</h3>
          <p className="text-xs text-studio-muted">Multi-faceted validation for narrative, wardrobe, and visual logic.</p>
        </Link>
      </div>
    </div>
  )
}
