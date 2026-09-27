import Link from 'next/link'
import { Plus, Clapperboard, ArrowRight } from 'lucide-react'

export default function ProjectsPage() {
  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-white">Production Projects</h1>
          <p className="text-xs text-studio-muted">Manage serialized drama series and active studio productions.</p>
        </div>
        <Link href="/projects/new" className="px-4 py-2 bg-studio-accent text-slate-950 font-bold text-xs rounded-lg flex items-center gap-2 hover:bg-cyan-400 transition">
          <Plus className="w-4 h-4" /> New Drama Project
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link href="/projects/prj-01" className="p-6 rounded-2xl glass-panel hover:border-studio-accent transition group">
          <div className="flex justify-between items-start mb-4">
            <div className="w-10 h-10 rounded-xl bg-studio-accent/10 border border-studio-accent/20 text-studio-accent flex items-center justify-center">
              <Clapperboard className="w-5 h-5" />
            </div>
            <span className="px-2.5 py-1 rounded-full text-[10px] font-bold bg-studio-violet/20 text-studio-violet border border-studio-violet/30">
              Autonomous Mode
            </span>
          </div>
          <h2 className="text-lg font-bold text-white group-hover:text-studio-accent transition mb-2">Shadows of Neo-London</h2>
          <p className="text-xs text-studio-muted mb-4 leading-relaxed">
            Cyberpunk mystery series following Detective Sarah Croft as she uncovers corporate conspiracies across Neo-London.
          </p>
          <div className="flex justify-between items-center text-[11px] text-studio-muted pt-4 border-t border-studio-border">
            <span>Season 1 • 12 Episodes</span>
            <span className="flex items-center gap-1 text-studio-accent font-semibold">Open Workspace <ArrowRight className="w-3.5 h-3.5" /></span>
          </div>
        </Link>
      </div>
    </div>
  )
}
