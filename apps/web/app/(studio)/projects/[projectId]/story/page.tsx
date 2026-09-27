export default function StoryPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Story & Narrative Hierarchy</h1>
      <div className="p-6 glass-panel rounded-2xl text-xs text-studio-muted space-y-4">
        <p className="text-white font-semibold text-sm">Series Narrative Tree: Series → Season 1 → Arc 2 → Episode 12</p>
        <p>Canonical story facts are stored independently in the Canon bounded context and resolved per episode.</p>
      </div>
    </div>
  )
}
