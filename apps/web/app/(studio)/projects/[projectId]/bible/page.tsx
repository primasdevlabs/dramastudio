"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { BookOpen, Save, Sparkles, CheckCircle2 } from "lucide-react";
import { api } from "@/lib/api/client";
import { SeriesBible } from "@/lib/api/types";

export default function SeriesBiblePage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const { data: bible, isLoading } = useQuery({
    queryKey: ["bible", projectId],
    queryFn: () => api.get<SeriesBible>(`/v1/projects/${projectId}/bible`).catch(() => null),
  });

  const [premise, setPremise] = useState("");
  const [genre, setGenre] = useState("Drama/Thriller");
  const [themes, setThemes] = useState("Betrayal, Power, Redemption");
  const [worldRules, setWorldRules] = useState("No magic, Grounded realism");
  const [narrativeRules, setNarrativeRules] = useState("Actions have permanent consequences");
  const [savedSuccess, setSavedSuccess] = useState(false);

  useEffect(() => {
    if (bible) {
      setPremise(bible.premise || "");
      setGenre(bible.genre || "Drama/Thriller");
      setThemes(bible.themes?.join(", ") || "");
      setWorldRules(bible.world_rules?.join("\n") || "");
      setNarrativeRules(bible.narrative_rules?.join("\n") || "");
    }
  }, [bible]);

  const saveMutation = useMutation({
    mutationFn: async () => {
      return api.post<SeriesBible>(`/v1/projects/${projectId}/bible`, {
        premise,
        genre,
        themes: themes.split(",").map((t) => t.trim()).filter(Boolean),
        world_rules: worldRules.split("\n").map((r) => r.trim()).filter(Boolean),
        narrative_rules: narrativeRules.split("\n").map((r) => r.trim()).filter(Boolean),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bible", projectId] });
      setSavedSuccess(true);
      setTimeout(() => setSavedSuccess(false), 3000);
    },
  });

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-center text-amber-400">
            <BookOpen className="w-5 h-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-xl font-bold text-white">Series Bible Workspace</h1>
              {bible && (
                <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 text-xs font-semibold border border-amber-500/20">
                  Version {bible.version}
                </span>
              )}
            </div>
            <p className="text-xs text-studio-muted">Canonical creative foundation for AI generation</p>
          </div>
        </div>

        <button
          onClick={() => saveMutation.mutate()}
          disabled={saveMutation.isPending}
          className="inline-flex items-center gap-2 bg-gradient-to-r from-amber-500 to-orange-600 hover:from-amber-400 hover:to-orange-500 text-black font-bold text-xs px-4 py-2.5 rounded-lg shadow-lg shadow-amber-500/20 transition-all"
        >
          <Save className="w-4 h-4" />
          {saveMutation.isPending ? "Saving Version..." : "Save Bible Version"}
        </button>
      </div>

      {savedSuccess && (
        <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 px-4 py-3 rounded-xl text-xs font-semibold flex items-center gap-2">
          <CheckCircle2 className="w-4 h-4" />
          Series Bible successfully saved and updated to next canonical version!
        </div>
      )}

      {/* Editor Form */}
      <div className="bg-studio-card border border-studio-border rounded-2xl p-6 space-y-6">
        <div>
          <label className="block text-xs font-bold text-studio-muted uppercase tracking-wider mb-2">
            Core Story Premise
          </label>
          <textarea
            rows={4}
            value={premise}
            onChange={(e) => setPremise(e.target.value)}
            placeholder="Write the central logline and core premise..."
            className="w-full bg-studio-panel border border-studio-border rounded-xl p-4 text-sm text-white focus:outline-none focus:border-amber-400 resize-none"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-xs font-bold text-studio-muted uppercase tracking-wider mb-2">
              Genre & Tone
            </label>
            <input
              type="text"
              value={genre}
              onChange={(e) => setGenre(e.target.value)}
              className="w-full bg-studio-panel border border-studio-border rounded-xl px-4 py-2.5 text-xs text-white focus:outline-none focus:border-amber-400"
            />
          </div>

          <div>
            <label className="block text-xs font-bold text-studio-muted uppercase tracking-wider mb-2">
              Central Themes (Comma Separated)
            </label>
            <input
              type="text"
              value={themes}
              onChange={(e) => setThemes(e.target.value)}
              className="w-full bg-studio-panel border border-studio-border rounded-xl px-4 py-2.5 text-xs text-white focus:outline-none focus:border-amber-400"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-xs font-bold text-studio-muted uppercase tracking-wider mb-2">
              World Rules (One per line)
            </label>
            <textarea
              rows={4}
              value={worldRules}
              onChange={(e) => setWorldRules(e.target.value)}
              placeholder="e.g. Grounded realism&#10;No supernatural events"
              className="w-full bg-studio-panel border border-studio-border rounded-xl p-3.5 text-xs text-white focus:outline-none focus:border-amber-400 resize-none font-mono"
            />
          </div>

          <div>
            <label className="block text-xs font-bold text-studio-muted uppercase tracking-wider mb-2">
              Narrative Rules (One per line)
            </label>
            <textarea
              rows={4}
              value={narrativeRules}
              onChange={(e) => setNarrativeRules(e.target.value)}
              placeholder="e.g. Actions have permanent consequences"
              className="w-full bg-studio-panel border border-studio-border rounded-xl p-3.5 text-xs text-white focus:outline-none focus:border-amber-400 resize-none font-mono"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
