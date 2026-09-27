"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Users, Plus, Shirt, UserCheck, Shield } from "lucide-react";
import { api } from "@/lib/api/client";
import { Character } from "@/lib/api/types";

export default function CharacterStudioPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [name, setName] = useState("");
  const [role, setRole] = useState("Protagonist");
  const [bio, setBio] = useState("");
  const [isOpen, setIsOpen] = useState(false);

  const { data: characters, isLoading } = useQuery({
    queryKey: ["characters", projectId],
    queryFn: async () => {
      const res = await api.get<{ characters: Character[] }>("/v1/characters");
      return res.characters || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      return api.post<Character>("/v1/characters", {
        project_id: projectId,
        name,
        role,
        bio,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["characters", projectId] });
      setIsOpen(false);
      setName("");
      setBio("");
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-pink-500/10 border border-pink-500/30 flex items-center justify-center text-pink-400">
            <Users className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Character Studio</h1>
            <p className="text-xs text-studio-muted">Manage persistent character identities, wardrobe, and voice profiles</p>
          </div>
        </div>

        <button
          onClick={() => setIsOpen(true)}
          className="inline-flex items-center gap-2 bg-gradient-to-r from-pink-500 to-rose-600 hover:from-pink-400 hover:to-rose-500 text-white font-bold text-xs px-4 py-2.5 rounded-lg shadow-lg shadow-pink-500/20 transition-all"
        >
          <Plus className="w-4 h-4" /> Add Character Profile
        </button>
      </div>

      {/* Modal */}
      {isOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-studio-card border border-studio-border rounded-xl w-full max-w-md p-6 shadow-2xl">
            <h2 className="text-lg font-bold text-white mb-4">New Character Profile</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Character Name</label>
                <input
                  type="text"
                  placeholder="e.g. Sarah Johnson"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-pink-400"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Role / Archetype</label>
                <input
                  type="text"
                  placeholder="e.g. Protagonist / Journalist"
                  value={role}
                  onChange={(e) => setRole(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-pink-400"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Biography & Appearance Notes</label>
                <textarea
                  rows={3}
                  placeholder="Character backstory and visual constraints..."
                  value={bio}
                  onChange={(e) => setBio(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-pink-400 resize-none"
                />
              </div>

              <div className="flex items-center justify-end gap-3 pt-4 border-t border-studio-border">
                <button
                  onClick={() => setIsOpen(false)}
                  className="px-4 py-2 text-xs font-semibold text-studio-muted hover:text-white"
                >
                  Cancel
                </button>
                <button
                  onClick={() => createMutation.mutate()}
                  disabled={createMutation.isPending || !name}
                  className="bg-pink-500 hover:bg-pink-400 text-white font-bold text-xs px-4 py-2 rounded-lg"
                >
                  {createMutation.isPending ? "Creating..." : "Save Character"}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isLoading ? (
          [1, 2].map((i) => <div key={i} className="h-48 bg-studio-card rounded-xl animate-pulse" />)
        ) : characters && characters.length > 0 ? (
          characters.map((c) => (
            <div
              key={c.id}
              className="bg-studio-card border border-studio-border rounded-xl p-5 space-y-3 hover:border-pink-500/40 transition-all"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-tr from-pink-500 to-purple-600 flex items-center justify-center font-bold text-white text-sm">
                    {c.name.substring(0, 2).toUpperCase()}
                  </div>
                  <div>
                    <h3 className="text-base font-bold text-white">{c.name}</h3>
                    <span className="text-xs text-pink-400 font-semibold">{c.role}</span>
                  </div>
                </div>
                <span className="px-2 py-0.5 rounded bg-pink-500/10 text-pink-400 text-xs font-mono border border-pink-500/20">
                  v{c.version || 1}
                </span>
              </div>

              <p className="text-xs text-studio-muted line-clamp-3 bg-studio-panel p-3 rounded-lg border border-studio-border/60">
                {c.bio || "No backstory notes configured."}
              </p>

              <div className="flex items-center justify-between text-xs pt-2 border-t border-studio-border/60">
                <span className="flex items-center gap-1 text-studio-muted">
                  <Shirt className="w-3.5 h-3.5 text-pink-400" /> Wardrobe Locked
                </span>
                <span className="flex items-center gap-1 text-emerald-400 font-medium">
                  <UserCheck className="w-3.5 h-3.5" /> Canon Reference
                </span>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full bg-studio-card border border-studio-border rounded-2xl p-8 text-center text-studio-muted text-xs">
            No characters created yet. Add a character to build your cast.
          </div>
        )}
      </div>
    </div>
  );
}
