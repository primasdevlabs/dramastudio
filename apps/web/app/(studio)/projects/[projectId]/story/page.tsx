export default function StoryPage() {
  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold text-white tracking-tight">Story & Narrative Hierarchy</h1>
      <div className="p-6 bg-studio-card border border-studio-border rounded-md text-xs text-studio-muted space-y-4">
        <p className="text-white font-semibold text-sm font-mono">Series Narrative Tree: Series → Season 1 → Arc 02 → Episode 12</p>
        <p className="leading-relaxed">Canonical story facts are stored independently in the Canon bounded context and resolved per episode through automated continuity checkers.</p>
      </div>
    </div>
  )
}
