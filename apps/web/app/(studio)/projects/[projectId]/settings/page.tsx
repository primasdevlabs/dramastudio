export default function SettingsPage() {
  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white tracking-tight">Studio Capability & Model Policy Settings</h1>
      
      <div className="p-6 bg-studio-card border border-studio-border rounded-md space-y-4">
        <h2 className="text-sm font-bold text-white uppercase tracking-wider font-mono">Capability Provider Mapping</h2>
        <div className="space-y-2 text-xs">
          <div className="p-3 bg-studio-panel rounded border border-studio-border flex justify-between items-center">
            <span className="font-semibold text-white">Scriptwriting & Dialogue</span>
            <span className="text-studio-accent font-mono">Model: Scriptwriter-v2 (Provider OpenAI)</span>
          </div>
          <div className="p-3 bg-studio-panel rounded border border-studio-border flex justify-between items-center">
            <span className="font-semibold text-white">Video Shot Generation (9:16)</span>
            <span className="text-studio-accent font-mono">Model: Wan 3.0 / Wan 2.1 (Provider Alibaba)</span>
          </div>
          <div className="p-3 bg-studio-panel rounded border border-studio-border flex justify-between items-center">
            <span className="font-semibold text-white">Voice & Dialogue Synthesis</span>
            <span className="text-emerald-400 font-mono">Model: ElevenLabs Multi-Voice (Provider ElevenLabs)</span>
          </div>
        </div>
      </div>
    </div>
  )
}
