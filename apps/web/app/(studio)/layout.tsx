'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Film, Layers, Users, Globe, Clapperboard, Sparkles, Activity, Settings, Radio, Sliders, ShieldCheck } from 'lucide-react'
import { useStudioStore } from '@/stores/studio-store'

export default function StudioLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const { productionMode, setProductionMode } = useStudioStore()

  const projectId = 'prj-01'

  const navItems = [
    { label: 'Project Dashboard', href: `/projects/${projectId}`, icon: Layers },
    { label: 'Story & Canon', href: `/projects/${projectId}/story`, icon: Sparkles },
    { label: 'Characters', href: `/projects/${projectId}/characters`, icon: Users },
    { label: 'World & Locations', href: `/projects/${projectId}/world`, icon: Globe },
    { label: 'Seasons & Episodes', href: `/projects/${projectId}/seasons`, icon: Clapperboard },
    { label: 'Production Activity', href: `/projects/${projectId}/production`, icon: Activity },
    { label: 'Continuity System', href: `/projects/${projectId}/continuity`, icon: ShieldCheck },
    { label: 'Automation & Agents', href: `/projects/${projectId}/automation`, icon: Radio },
    { label: 'Publishing', href: `/projects/${projectId}/publishing`, icon: Sliders },
    { label: 'Studio Settings', href: `/projects/${projectId}/settings`, icon: Settings },
  ]

  return (
    <div className="min-h-screen flex bg-studio-bg text-studio-text">
      {/* Sidebar */}
      <aside className="w-64 border-r border-studio-border bg-studio-card flex flex-col justify-between p-4 shrink-0">
        <div>
          <div className="flex items-center gap-3 px-2 py-3 mb-6 border-b border-studio-border">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-studio-accent to-studio-violet flex items-center justify-center glow-cyan">
              <Film className="w-4 h-4 text-white" />
            </div>
            <div>
              <div className="font-bold text-sm text-white leading-tight">DramaStudio</div>
              <div className="text-[10px] text-studio-accent font-semibold">Virtual Studio Engine</div>
            </div>
          </div>

          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = pathname.startsWith(item.href)
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-xs font-medium transition ${
                    isActive
                      ? 'bg-studio-panel text-studio-accent font-semibold border-l-2 border-studio-accent'
                      : 'text-studio-muted hover:text-white hover:bg-studio-panel/50'
                  }`}
                >
                  <Icon className="w-4 h-4" />
                  {item.label}
                </Link>
              )
            })}
          </nav>
        </div>

        {/* Mode Selector */}
        <div className="p-3 rounded-xl bg-studio-panel border border-studio-border space-y-2">
          <div className="text-[10px] uppercase font-bold tracking-wider text-studio-muted">Studio Operating Mode</div>
          <div className="grid grid-cols-2 gap-1 bg-studio-card p-1 rounded-lg">
            <button
              onClick={() => setProductionMode('monitored')}
              className={`py-1.5 text-[11px] font-semibold rounded-md transition ${
                productionMode === 'monitored' ? 'bg-studio-accent text-slate-950 shadow' : 'text-studio-muted hover:text-white'
              }`}
            >
              Monitored
            </button>
            <button
              onClick={() => setProductionMode('autonomous')}
              className={`py-1.5 text-[11px] font-semibold rounded-md transition ${
                productionMode === 'autonomous' ? 'bg-studio-violet text-white shadow' : 'text-studio-muted hover:text-white'
              }`}
            >
              Autonomous
            </button>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <header className="h-14 border-b border-studio-border bg-studio-card px-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="text-xs text-studio-muted">Project:</span>
            <span className="text-xs font-bold text-white bg-studio-panel px-2.5 py-1 rounded border border-studio-border">
              Shadows of Neo-London (Season 1)
            </span>
          </div>

          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-[11px] font-medium">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              Lead Director Active
            </div>
          </div>
        </header>

        {/* Page Viewport */}
        <main className="flex-1 p-6 overflow-y-auto">
          {children}
        </main>
      </div>
    </div>
  )
}
