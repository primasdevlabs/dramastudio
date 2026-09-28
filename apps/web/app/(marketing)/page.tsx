import Link from 'next/link'
import {
  Film, Bot, Layers, Cpu, ArrowRight, Clapperboard, BookOpen,
  GitBranch, ShieldCheck, MonitorPlay, Share2, Check, Sparkles,
  PenLine, Eye, Rocket,
} from 'lucide-react'

const PIPELINE_STAGES = [
  { label: 'Script', icon: PenLine },
  { label: 'Storyboard', icon: Clapperboard },
  { label: 'Shots', icon: Film },
  { label: 'Render', icon: MonitorPlay },
  { label: 'Publish', icon: Share2 },
]

const FEATURES = [
  {
    icon: Bot,
    accent: 'text-studio-accent',
    variant: '',
    title: 'Lead Director Console',
    body: 'An orchestration layer of specialized agents — story, visual, audio, assembly, and continuity — executing under a single director loop. Every decision is logged, inspectable, and reversible.',
  },
  {
    icon: BookOpen,
    accent: 'text-studio-amber',
    variant: 'grain-card--amber',
    title: 'Persistent Story Canon',
    body: 'A fact knowledge graph tracks every established truth — character backstory, world rules, plot threads — with full version history and retcon support across 150-episode arcs.',
  },
  {
    icon: ShieldCheck,
    accent: 'text-studio-green',
    variant: 'grain-card--green',
    title: 'Continuity Engine',
    body: 'Automated multi-faceted checkers validate each episode against canon before release. Contradictions surface as issues with evidence, cause, and resolution tracking.',
  },
  {
    icon: Cpu,
    accent: 'text-studio-amber',
    variant: 'grain-card--amber',
    title: 'Capability Model Registry',
    body: 'Provider-agnostic routing across sixteen adapters — OpenAI, Anthropic, ElevenLabs, Suno, Wan, and more. Swap models per capability without touching a pipeline.',
  },
  {
    icon: GitBranch,
    accent: 'text-studio-accent',
    variant: '',
    title: 'Approval Gates',
    body: 'Monitored mode pauses the run at script, storyboard, and final cut. Review the exact artifacts the pipeline produced, decide, and the run resumes from that checkpoint.',
  },
  {
    icon: Share2,
    accent: 'text-studio-red',
    variant: 'grain-card--red',
    title: 'Render & Publish',
    body: 'FFmpeg assembly builds the master timeline from approved shots, then scheduled publications ship to connected channels with per-platform idempotent delivery.',
  },
]

const PLANS = [
  {
    name: 'Monitored Studio',
    price: '$99',
    tagline: 'Human approval at every major checkpoint.',
    features: [
      'Script, storyboard, and final-cut gates',
      'Full canon & continuity validation',
      'Single-channel publishing',
      'Bring your own provider keys',
    ],
    cta: 'Start Monitored Production',
    featured: false,
  },
  {
    name: 'Autonomous Studio',
    price: '$299',
    tagline: 'The Lead Director runs the loop end to end.',
    features: [
      'Autonomous director loop execution',
      'FFmpeg auto-assembly pipeline',
      'Multi-channel scheduled publishing',
      'Priority provider routing policies',
    ],
    cta: 'Start Autonomous Studio',
    featured: true,
  },
]

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-studio-bg text-studio-text">
      {/* ── Nav ── */}
      <header className="fixed top-0 inset-x-0 z-50 border-b border-studio-border/60 bg-studio-bg/80 backdrop-blur-md">
        <div className="max-w-7xl mx-auto flex justify-between items-center px-6 py-4">
          <Link href="/" className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-md bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
              <Film className="w-4.5 h-4.5" />
            </div>
            <span className="font-bold text-lg tracking-tight text-white">DramaStudio</span>
          </Link>
          <nav className="hidden md:flex items-center gap-8">
            <a href="#features" className="text-sm text-studio-muted hover:text-white transition">Features</a>
            <a href="#about" className="text-sm text-studio-muted hover:text-white transition">About</a>
            <a href="#pricing" className="text-sm text-studio-muted hover:text-white transition">Pricing</a>
          </nav>
          <div className="flex items-center gap-3">
            <Link href="/login" className="text-sm text-studio-muted hover:text-white transition hidden sm:block">
              Sign in
            </Link>
            <Link
              href="/register"
              className="px-4 py-2 rounded bg-studio-accent text-white font-medium text-sm hover:bg-studio-accent-dark transition flex items-center gap-2"
            >
              Enter Studio <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </div>
      </header>

      {/* ── Hero ── */}
      <section className="relative pt-36 pb-24 px-6 overflow-hidden">
        <div className="glow-field w-[600px] h-[400px] bg-studio-accent/15 -top-40 left-1/2 -translate-x-1/2" />
        <div className="glow-field w-[300px] h-[300px] bg-studio-amber/10 top-40 -left-32" />

        <div className="relative max-w-5xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-studio-panel/80 border border-studio-border text-studio-accent text-xs font-mono uppercase tracking-wider mb-8">
            <Sparkles className="w-3.5 h-3.5" /> Serialized Drama Production Workstation
          </div>
          <h1 className="text-5xl md:text-7xl font-bold tracking-tight text-white mb-6 leading-[1.05]">
            Automated Production.<br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-studio-accent via-studio-amber to-studio-accent">
              Persistent Story Canon.
            </span>
          </h1>
          <p className="text-lg md:text-xl text-studio-muted max-w-2xl mx-auto mb-10 leading-relaxed">
            Orchestrate specialized model backends through a Lead Director console. Turn story
            premises into 150-episode serialized dramas with continuity that never drifts.
          </p>
          <div className="flex flex-wrap justify-center gap-4 mb-20">
            <Link
              href="/register"
              className="px-8 py-4 rounded bg-studio-accent text-white font-bold hover:bg-studio-accent-dark transition flex items-center gap-3"
            >
              Open Studio Console <ArrowRight className="w-5 h-5" />
            </Link>
            <a
              href="#features"
              className="px-8 py-4 rounded bg-studio-panel/80 border border-studio-border text-white font-medium hover:border-studio-accent transition"
            >
              See how it works
            </a>
          </div>

          {/* Pipeline visual */}
          <div className="grain-card rounded-lg max-w-3xl mx-auto p-6 md:p-8">
            <div className="flex items-center justify-between gap-2">
              {PIPELINE_STAGES.map((stage, i) => (
                <div key={stage.label} className="flex items-center gap-2 flex-1 last:flex-none">
                  <div className="flex flex-col items-center gap-2 flex-1">
                    <div className="w-11 h-11 rounded-md bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
                      <stage.icon className="w-5 h-5" />
                    </div>
                    <span className="text-[10px] font-mono uppercase tracking-wider text-studio-muted">
                      {stage.label}
                    </span>
                  </div>
                  {i < PIPELINE_STAGES.length - 1 && (
                    <div className="hidden sm:block h-px flex-1 bg-gradient-to-r from-studio-border to-studio-accent/40 mb-6" />
                  )}
                </div>
              ))}
            </div>
            <p className="text-center text-xs text-studio-muted mt-6 font-mono">
              One director loop · five production stages · zero canon drift
            </p>
          </div>
        </div>
      </section>

      {/* ── Stats strip ── */}
      <section className="border-y border-studio-border/60 bg-studio-card/40">
        <div className="max-w-6xl mx-auto grid grid-cols-2 md:grid-cols-4 divide-x divide-studio-border/60">
          {[
            ['150', 'episode arcs per series'],
            ['16', 'provider adapters routed'],
            ['5', 'specialist agent roles'],
            ['2', 'production modes'],
          ].map(([num, label]) => (
            <div key={label} className="px-6 py-8 text-center">
              <div className="text-3xl font-bold text-white font-mono">{num}</div>
              <div className="text-xs text-studio-muted mt-1">{label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── About ── */}
      <section id="about" className="relative py-24 px-6 scroll-mt-20 overflow-hidden">
        <div className="glow-field w-[400px] h-[300px] bg-studio-accent/10 top-20 -right-40" />
        <div className="relative max-w-6xl mx-auto">
          <div className="max-w-2xl mb-14">
            <p className="text-xs font-mono uppercase tracking-widest text-studio-accent mb-3">About</p>
            <h2 className="text-3xl md:text-4xl font-bold text-white tracking-tight mb-4">
              Built like a real studio, not a prompt chain.
            </h2>
            <p className="text-studio-muted leading-relaxed">
              Generating episodes is easy. Generating a hundred episodes where character twelve
              still fears the water and the season-two prophecy actually pays off — that takes a
              system. DramaStudio gives every episode the same canon, the same checkers, and the
              same director, so volume never costs you coherence.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
            {[
              {
                icon: PenLine,
                step: '01',
                title: 'Write',
                body: 'Premises become season outlines, arcs, and episodes. Canon facts are extracted and versioned as the story is written — not bolted on after.',
              },
              {
                icon: Eye,
                step: '02',
                title: 'Direct',
                body: 'The Lead Director sequences specialist agents through script, storyboard, shots, and voice — pausing at approval gates when you want eyes on the work.',
              },
              {
                icon: Rocket,
                step: '03',
                title: 'Ship',
                body: 'Continuity-verified episodes are assembled, rendered, and scheduled to your channels. Publishing is idempotent — a retry never double-posts.',
              },
            ].map((s) => (
              <div key={s.step} className="grain-card rounded-lg p-6">
                <div className="flex items-center justify-between mb-5">
                  <s.icon className="w-6 h-6 text-studio-accent" />
                  <span className="text-xs font-mono text-studio-muted">{s.step}</span>
                </div>
                <h3 className="font-bold text-white mb-2">{s.title}</h3>
                <p className="text-sm text-studio-muted leading-relaxed">{s.body}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section id="features" className="relative py-24 px-6 border-t border-studio-border/60 scroll-mt-20">
        <div className="max-w-6xl mx-auto">
          <div className="text-center max-w-2xl mx-auto mb-14">
            <p className="text-xs font-mono uppercase tracking-widest text-studio-accent mb-3">Features</p>
            <h2 className="text-3xl md:text-4xl font-bold text-white tracking-tight mb-4">
              The whole pipeline, under one roof
            </h2>
            <p className="text-studio-muted">
              Every layer a serialized production needs — orchestration, memory, validation,
              rendering, and distribution — sharing one canon.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {FEATURES.map((f) => (
              <div
                key={f.title}
                className={`grain-card ${f.variant} rounded-lg p-6 hover:border-studio-accent/60 transition-colors`}
              >
                <f.icon className={`w-8 h-8 ${f.accent} mb-4`} />
                <h3 className="font-bold text-white mb-2">{f.title}</h3>
                <p className="text-sm text-studio-muted leading-relaxed">{f.body}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Pricing ── */}
      <section id="pricing" className="relative py-24 px-6 border-t border-studio-border/60 scroll-mt-20 overflow-hidden">
        <div className="glow-field w-[500px] h-[300px] bg-studio-accent/10 -bottom-20 left-1/3" />
        <div className="relative max-w-4xl mx-auto">
          <div className="text-center max-w-2xl mx-auto mb-14">
            <p className="text-xs font-mono uppercase tracking-widest text-studio-accent mb-3">Pricing</p>
            <h2 className="text-3xl md:text-4xl font-bold text-white tracking-tight mb-4">
              Production studio plans
            </h2>
            <p className="text-studio-muted">
              Model provider capacity, director loop orchestration, and multi-track rendering.
              Bring your own provider keys — spend stays visible and under your control.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {PLANS.map((plan) => (
              <div
                key={plan.name}
                className={`grain-card rounded-lg p-8 flex flex-col ${
                  plan.featured ? 'border-studio-accent' : ''
                }`}
              >
                {plan.featured && (
                  <span className="self-start text-[10px] font-mono uppercase tracking-widest text-studio-accent border border-studio-accent/40 rounded-full px-2.5 py-1 mb-4">
                    Recommended
                  </span>
                )}
                <h3 className="text-xl font-bold text-white">{plan.name}</h3>
                <p className="text-3xl font-bold text-white font-mono mt-3">
                  {plan.price}
                  <span className="text-sm text-studio-muted font-sans font-normal"> / month</span>
                </p>
                <p className="text-sm text-studio-muted mt-2 mb-6">{plan.tagline}</p>
                <ul className="space-y-3 mb-8 flex-1">
                  {plan.features.map((feat) => (
                    <li key={feat} className="flex items-start gap-2.5 text-sm text-studio-muted">
                      <Check className="w-4 h-4 text-studio-green mt-0.5 shrink-0" />
                      {feat}
                    </li>
                  ))}
                </ul>
                <Link
                  href="/register"
                  className={`block text-center py-2.5 rounded font-medium text-sm transition ${
                    plan.featured
                      ? 'bg-studio-accent text-white hover:bg-studio-accent-dark'
                      : 'bg-studio-panel border border-studio-border text-white hover:border-studio-accent'
                  }`}
                >
                  {plan.cta}
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── CTA band ── */}
      <section className="relative py-24 px-6 border-t border-studio-border/60 overflow-hidden">
        <div className="glow-field w-[700px] h-[300px] bg-studio-accent/15 top-0 left-1/2 -translate-x-1/2" />
        <div className="relative grain-card rounded-xl max-w-4xl mx-auto p-10 md:p-14 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-white tracking-tight mb-4">
            Your writers' room never sleeps
          </h2>
          <p className="text-studio-muted max-w-xl mx-auto mb-8">
            Register a studio, drop in a premise, and watch the first episode work its way
            through the pipeline — script to rendered master.
          </p>
          <Link
            href="/register"
            className="inline-flex items-center gap-3 px-8 py-4 rounded bg-studio-accent text-white font-bold hover:bg-studio-accent-dark transition"
          >
            Create your studio <ArrowRight className="w-5 h-5" />
          </Link>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="border-t border-studio-border/60 bg-studio-card/40">
        <div className="max-w-7xl mx-auto px-6 py-12 grid grid-cols-2 md:grid-cols-4 gap-8">
          <div className="col-span-2 md:col-span-1">
            <div className="flex items-center gap-2.5 mb-4">
              <div className="w-8 h-8 rounded-md bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
                <Film className="w-4 h-4" />
              </div>
              <span className="font-bold text-white">DramaStudio</span>
            </div>
            <p className="text-xs text-studio-muted leading-relaxed">
              The serialized drama production workstation. Automated pipeline, persistent canon.
            </p>
          </div>
          <div>
            <h4 className="text-xs font-mono uppercase tracking-widest text-white mb-4">Product</h4>
            <ul className="space-y-2.5 text-sm text-studio-muted">
              <li><a href="#features" className="hover:text-white transition">Features</a></li>
              <li><a href="#pricing" className="hover:text-white transition">Pricing</a></li>
              <li><Link href="/projects" className="hover:text-white transition">Studio Console</Link></li>
            </ul>
          </div>
          <div>
            <h4 className="text-xs font-mono uppercase tracking-widest text-white mb-4">Studio</h4>
            <ul className="space-y-2.5 text-sm text-studio-muted">
              <li><a href="#about" className="hover:text-white transition">About</a></li>
              <li><Link href="/register" className="hover:text-white transition">Create account</Link></li>
              <li><Link href="/login" className="hover:text-white transition">Sign in</Link></li>
            </ul>
          </div>
          <div>
            <h4 className="text-xs font-mono uppercase tracking-widest text-white mb-4">Legal</h4>
            <ul className="space-y-2.5 text-sm text-studio-muted">
              <li><span className="cursor-default">Terms of Service</span></li>
              <li><span className="cursor-default">Privacy Policy</span></li>
            </ul>
          </div>
        </div>
        <div className="border-t border-studio-border/60">
          <div className="max-w-7xl mx-auto px-6 py-5 flex flex-col sm:flex-row justify-between items-center gap-2 text-xs text-studio-muted">
            <span>Copyright © 2026 DramaStudio. All rights reserved.</span>
            <span className="font-mono">script → storyboard → shots → render → publish</span>
          </div>
        </div>
      </footer>
    </div>
  )
}
