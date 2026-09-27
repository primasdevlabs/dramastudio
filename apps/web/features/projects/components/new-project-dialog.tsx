"use client";

import { useState } from "react";
import { Plus, Clapperboard, Sparkles } from "lucide-react";
import { api } from "@/lib/api/client";
import { Project } from "@/lib/api/types";

export function NewProjectDialog({ onCreated }: { onCreated?: (p: Project) => void }) {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [genre, setGenre] = useState("Drama/Thriller");
  const [language, setLanguage] = useState("en");
  const [mode, setMode] = useState<"monitored" | "autonomous">("monitored");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setIsSubmitting(true);
    try {
      const proj = await api.post<Project>("/v1/projects", {
        name,
        description,
        genre,
        language,
        mode,
      });
      setIsOpen(false);
      setName("");
      setDescription("");
      if (onCreated) onCreated(proj);
    } catch (err) {
      console.error("Failed to create project", err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="inline-flex items-center gap-2 bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white font-semibold text-sm px-4 py-2 rounded-lg shadow-lg shadow-cyan-500/20 transition-all"
      >
        <Plus className="w-4 h-4" /> New Production
      </button>

      {isOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-studio-card border border-studio-border rounded-xl w-full max-w-md p-6 shadow-2xl">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-lg bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400">
                <Sparkles className="w-5 h-5" />
              </div>
              <div>
                <h2 className="text-lg font-bold text-white">Create Drama Production</h2>
                <p className="text-xs text-studio-muted">Set up your new AI drama series workspace</p>
              </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Production Title</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. The Last Promise"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-cyan-400"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Premise / Description</label>
                <textarea
                  rows={3}
                  placeholder="Brief story premise..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-cyan-400 resize-none"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-studio-muted mb-1">Genre</label>
                  <select
                    value={genre}
                    onChange={(e) => setGenre(e.target.value)}
                    className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-cyan-400"
                  >
                    <option value="Drama/Thriller">Drama / Thriller</option>
                    <option value="Romance">Romance</option>
                    <option value="Sci-Fi">Sci-Fi</option>
                    <option value="Action">Action</option>
                    <option value="Mystery">Mystery</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-studio-muted mb-1">Production Mode</label>
                  <select
                    value={mode}
                    onChange={(e) => setMode(e.target.value as any)}
                    className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-cyan-400"
                  >
                    <option value="monitored">Monitored (Human Approval)</option>
                    <option value="autonomous">Autonomous (Lead Director)</option>
                  </select>
                </div>
              </div>

              <div className="flex items-center justify-end gap-3 pt-4 border-t border-studio-border">
                <button
                  type="button"
                  onClick={() => setIsOpen(false)}
                  className="px-4 py-2 text-xs font-semibold text-studio-muted hover:text-white transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="bg-cyan-500 hover:bg-cyan-400 text-black font-bold text-xs px-4 py-2 rounded-lg transition-all"
                >
                  {isSubmitting ? "Creating..." : "Initialize Production"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
