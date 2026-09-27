import Link from 'next/link'

export default function NewProjectPage() {
  return (
    <div className="max-w-2xl mx-auto p-8 bg-studio-card border border-studio-border rounded-md">
      <h1 className="text-xl font-bold text-white mb-2 tracking-tight">Initialize Drama Production</h1>
      <p className="text-xs text-studio-muted mb-6">Enter a core story logline or prompt. The Lead Director Agent will generate the Series Bible.</p>

      <form className="space-y-4">
        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Series Title</label>
          <input type="text" placeholder="e.g. Echoes of Tomorrow" className="w-full bg-studio-panel border border-studio-border rounded px-4 py-2.5 text-xs text-white focus:outline-none focus:border-studio-accent" />
        </div>
        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Initial Story Concept / Logline</label>
          <textarea rows={4} placeholder="Describe the central premise, protagonist, and core conflict..." className="w-full bg-studio-panel border border-studio-border rounded p-4 text-xs text-white focus:outline-none focus:border-studio-accent" />
        </div>
        <div className="flex justify-end gap-3 pt-4">
          <Link href="/projects" className="px-4 py-2 text-xs text-studio-muted hover:text-white">Cancel</Link>
          <Link href="/projects/prj-01" className="px-5 py-2 bg-studio-accent hover:bg-studio-accent-dark text-white font-medium text-xs rounded transition-colors">Generate Series Bible</Link>
        </div>
      </form>
    </div>
  )
}
