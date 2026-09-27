export default function SettingsPage() {
  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white">Studio Capability & Model Policy Settings</h1>
      
      <div className="p-6 glass-panel rounded-2xl space-y-4">
        <h2 className="text-sm font-bold text-white">Capability Provider Mapping</h2>
        <div className="space-y-2 text-xs">
          <div className="p-3 bg-studio-bg rounded-lg border border-studio-border flex justify-between items-center">
            <span>ScriptWriting</span>
            <span className="text-studio-accent font-semibold">Model: Script-v2 (Provider A)</span>
          </div>
          <div className="p-3 bg-studio-bg rounded-lg border border-studio-border flex justify-between items-center">
            <span>Video Generation</span>
            <span className="text-studio-violet font-semibold">Model: Wan 2.1 (External Provider C)</span>
          </div>
          <div className="p-3 bg-studio-bg rounded-lg border border-studio-border flex justify-between items-center">
            <span>Voice Generation</span>
            <span className="text-emerald-400 font-semibold">Model: TTS-HD (Provider D)</span>
          </div>
        </div>
      </div>
    </div>
  )
}
