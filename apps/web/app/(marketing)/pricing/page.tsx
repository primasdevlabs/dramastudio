import Link from 'next/link'

export default function PricingPage() {
  return (
    <div className="min-h-screen bg-studio-bg p-8 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold text-white mb-4">Production Studio Plans</h1>
      <p className="text-studio-muted mb-8">Model provider API capacity, autonomous director loop orchestration, and multi-track rendering.</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="p-6 rounded-md bg-studio-card border border-studio-border">
          <h2 className="text-xl font-bold text-white mb-2">Monitored Mode Studio</h2>
          <p className="text-2xl font-bold text-studio-accent mb-4 font-mono">$99 <span className="text-xs text-studio-muted">/ month</span></p>
          <p className="text-sm text-studio-muted mb-6">Human approval required at major production checkpoints (script, storyboard, final cut).</p>
          <Link href="/projects" className="block text-center py-2.5 rounded bg-studio-panel border border-studio-border text-white font-medium text-sm hover:border-studio-accent transition-colors">Start Monitored Production</Link>
        </div>
        <div className="p-6 rounded-md bg-studio-card border border-studio-accent">
          <h2 className="text-xl font-bold text-white mb-2">Autonomous Studio</h2>
          <p className="text-2xl font-bold text-studio-accent mb-4 font-mono">$299 <span className="text-xs text-studio-muted">/ month</span></p>
          <p className="text-sm text-studio-muted mb-6">Lead Director agent autonomous loop execution, FFmpeg auto-assembly & multi-channel publishing.</p>
          <Link href="/projects" className="block text-center py-2.5 rounded bg-studio-accent text-white font-medium text-sm hover:bg-studio-accent-dark transition-colors">Start Autonomous Studio</Link>
        </div>
      </div>
    </div>
  )
}
