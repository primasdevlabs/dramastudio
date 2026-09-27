import Link from 'next/link'

export default function PricingPage() {
  return (
    <div className="min-h-screen bg-studio-bg p-8 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold text-white mb-4">Production Plans</h1>
      <p className="text-studio-muted mb-8">Pay for model provider API capacity and studio orchestration.</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="p-6 rounded-2xl glass-panel border-studio-accent">
          <h2 className="text-xl font-bold text-white mb-2">Monitored Mode Studio</h2>
          <p className="text-2xl font-bold text-studio-accent mb-4">$99 <span className="text-xs text-studio-muted">/ month</span></p>
          <p className="text-sm text-studio-muted mb-6">Human approval required at major production checkpoints.</p>
          <Link href="/projects" className="block text-center py-3 rounded-lg bg-studio-accent text-slate-950 font-bold">Start Production</Link>
        </div>
        <div className="p-6 rounded-2xl glass-panel border-studio-violet">
          <h2 className="text-xl font-bold text-white mb-2">Autonomous Studio</h2>
          <p className="text-2xl font-bold text-studio-violet mb-4">$299 <span className="text-xs text-studio-muted">/ month</span></p>
          <p className="text-sm text-studio-muted mb-6">Lead Director agent autonomous execution & auto-publishing.</p>
          <Link href="/projects" className="block text-center py-3 rounded-lg bg-studio-violet text-white font-bold">Start Autonomous Studio</Link>
        </div>
      </div>
    </div>
  )
}
