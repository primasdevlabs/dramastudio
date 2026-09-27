"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ShieldCheck, Plus, CheckCircle2, AlertTriangle, XCircle } from "lucide-react";
import { api } from "@/lib/api/client";
import { StoryFact } from "@/lib/api/types";

export default function CanonPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [subject, setSubject] = useState("");
  const [predicate, setPredicate] = useState("owns");
  const [object, setObject] = useState("");
  const [isOpen, setIsOpen] = useState(false);

  const { data: facts, isLoading } = useQuery({
    queryKey: ["canon-facts", projectId],
    queryFn: async () => {
      const res = await api.get<{ facts: StoryFact[] }>(`/v1/canon/facts?project_id=${projectId}`);
      return res.facts || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      return api.post<StoryFact>("/v1/canon/facts", {
        subject,
        predicate,
        object,
        introduced: "Episode 1",
        valid_from: "Episode 1",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["canon-facts", projectId] });
      setIsOpen(false);
      setSubject("");
      setObject("");
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400">
            <ShieldCheck className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Canon Context & Story Facts</h1>
            <p className="text-xs text-studio-muted">Authoritative source of story truth and character knowledge state</p>
          </div>
        </div>

        <button
          onClick={() => setIsOpen(true)}
          className="inline-flex items-center gap-2 bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white font-bold text-xs px-4 py-2.5 rounded-lg shadow-lg shadow-cyan-500/20 transition-all"
        >
          <Plus className="w-4 h-4" /> Add Canonical Fact
        </button>
      </div>

      {/* Modal */}
      {isOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-studio-card border border-studio-border rounded-xl w-full max-w-md p-6 shadow-2xl">
            <h2 className="text-lg font-bold text-white mb-4">Establish Story Fact</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Subject (Character / Entity)</label>
                <input
                  type="text"
                  placeholder="e.g. Sarah"
                  value={subject}
                  onChange={(e) => setSubject(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-cyan-400"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Predicate (Relation / Fact)</label>
                <input
                  type="text"
                  placeholder="e.g. discovered_secret"
                  value={predicate}
                  onChange={(e) => setPredicate(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-cyan-400"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Object (Target / Information)</label>
                <input
                  type="text"
                  placeholder="e.g. confidential_file"
                  value={object}
                  onChange={(e) => setObject(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-cyan-400"
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
                  disabled={createMutation.isPending || !subject || !object}
                  className="bg-cyan-500 hover:bg-cyan-400 text-black font-bold text-xs px-4 py-2 rounded-lg"
                >
                  {createMutation.isPending ? "Saving..." : "Establish Fact"}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Facts Table */}
      <div className="bg-studio-card border border-studio-border rounded-2xl overflow-hidden shadow-xl">
        <table className="w-full text-left text-xs text-studio-text">
          <thead className="bg-studio-panel border-b border-studio-border text-studio-muted uppercase tracking-wider font-semibold">
            <tr>
              <th className="p-4">Subject</th>
              <th className="p-4">Predicate</th>
              <th className="p-4">Object</th>
              <th className="p-4">Introduced In</th>
              <th className="p-4">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-studio-border/60">
            {isLoading ? (
              <tr>
                <td colSpan={5} className="p-8 text-center text-studio-muted italic">
                  Loading canonical story facts...
                </td>
              </tr>
            ) : facts && facts.length > 0 ? (
              facts.map((fact) => (
                <tr key={fact.id} className="hover:bg-studio-panel/50 transition-colors">
                  <td className="p-4 font-bold text-white">{fact.subject}</td>
                  <td className="p-4 font-mono text-cyan-400">{fact.predicate}</td>
                  <td className="p-4 font-semibold text-studio-text">{fact.object}</td>
                  <td className="p-4 text-studio-muted">{fact.introduced || "Episode 1"}</td>
                  <td className="p-4">
                    {fact.status === "canonical" ? (
                      <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 font-semibold">
                        <CheckCircle2 className="w-3.5 h-3.5" /> Canonical
                      </span>
                    ) : fact.status === "disputed" ? (
                      <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20 font-semibold">
                        <AlertTriangle className="w-3.5 h-3.5" /> Disputed
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded bg-rose-500/10 text-rose-400 border border-rose-500/20 font-semibold">
                        <XCircle className="w-3.5 h-3.5" /> Retconned
                      </span>
                    )}
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={5} className="p-8 text-center text-studio-muted italic">
                  No canonical facts established yet.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
