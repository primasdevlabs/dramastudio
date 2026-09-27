import Link from 'next/link'

export default function NewProjectPage() {
  return (
    <div className="max-w-2xl mx-auto p-8 glass-panel rounded-2xl">
      <h1 className="text-xl font-bold text-white mb-2">Initialize Drama Project</h1>
      <p className="text-xs text-studio-muted mb-6">Enter a core story logline or prompt. The Lead Director Agent will generate the Series Bible.</p>

      <form className="space-y-4">
        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Series Title</label>
          <input type="text" placeholder="e.g. Echoes of Tomorrow" className="w-full bg-studio-bg border border-studio-border rounded-lg px-4 py-2.5 text-xs text-white focus:outline-none focus:border-studio-accent" />
        </div>
        <div>
          <label className="block text-xs font-semibold text-studio-muted mb-1">Initial Story Concept / Logline</label>
          <textarea rows={4} placeholder="Describe the central premise, protagonist, and core conflict..." className="w-full bg-studio-bg border border-studio-border rounded-lg p-4 text-xs text-white focus:outline-none focus:border-studio-accent" />
        </div>
        <div className="flex justify-end gap-3 pt-4">
          <Link href="/projects" className="px-4 py-2 text-xs text-studio-muted hover:text-white">Cancel</Link>
          <Link href="/projects/prj-01" className="px-5 py-2 bg-studio-accent text-slate-950 font-bold text-xs rounded-lg">Generate Series Bible</Link>
        </div>
      </form>
    </div>
  )
}
