"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Globe, Plus, Sun, Moon, CloudRain, MapPin } from "lucide-react";
import { api } from "@/lib/api/client";
import { Location } from "@/lib/api/types";

export default function WorldStudioPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [type, setType] = useState("Interior");
  const [isOpen, setIsOpen] = useState(false);

  const { data: locations, isLoading } = useQuery({
    queryKey: ["locations", projectId],
    queryFn: async () => {
      const res = await api.get<{ locations: Location[] }>(`/v1/world/locations?project_id=${projectId}`);
      return res.locations || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      return api.post<Location>("/v1/world/locations", {
        project_id: projectId,
        name,
        description,
        type,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["locations", projectId] });
      setIsOpen(false);
      setName("");
      setDescription("");
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center text-emerald-400">
            <Globe className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">World Studio</h1>
            <p className="text-xs text-studio-muted">Manage persistent production locations, environments, and lighting variants</p>
          </div>
        </div>

        <button
          onClick={() => setIsOpen(true)}
          className="inline-flex items-center gap-2 bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-400 hover:to-teal-500 text-black font-bold text-xs px-4 py-2.5 rounded-lg shadow-lg shadow-emerald-500/20 transition-all"
        >
          <Plus className="w-4 h-4" /> Add Production Location
        </button>
      </div>

      {/* Modal */}
      {isOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-studio-card border border-studio-border rounded-xl w-full max-w-md p-6 shadow-2xl">
            <h2 className="text-lg font-bold text-white mb-4">New Production Location</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Location Name</label>
                <input
                  type="text"
                  placeholder="e.g. Downtown Apartment"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-400"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Type</label>
                <select
                  value={type}
                  onChange={(e) => setType(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-emerald-400"
                >
                  <option value="Interior">Interior</option>
                  <option value="Exterior">Exterior</option>
                  <option value="Studio Set">Studio Set</option>
                </select>
              </div>

              <div>
                <label className="block text-xs font-semibold text-studio-muted mb-1">Description & Visual Rules</label>
                <textarea
                  rows={3}
                  placeholder="Describe lighting, decor, and spatial layout..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="w-full bg-studio-panel border border-studio-border rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-emerald-400 resize-none"
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
                  className="bg-emerald-500 hover:bg-emerald-400 text-black font-bold text-xs px-4 py-2 rounded-lg"
                >
                  {createMutation.isPending ? "Saving..." : "Save Location"}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isLoading ? (
          [1, 2].map((i) => <div key={i} className="h-44 bg-studio-card rounded-xl animate-pulse" />)
        ) : locations && locations.length > 0 ? (
          locations.map((loc) => (
            <div
              key={loc.id}
              className="bg-studio-card border border-studio-border rounded-xl p-5 space-y-3 hover:border-emerald-500/40 transition-all"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <MapPin className="w-4 h-4 text-emerald-400" />
                  <h3 className="text-base font-bold text-white">{loc.name}</h3>
                </div>
                <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 text-xs font-semibold border border-emerald-500/20">
                  {loc.type}
                </span>
              </div>

              <p className="text-xs text-studio-muted line-clamp-2 bg-studio-panel p-3 rounded-lg border border-studio-border/60">
                {loc.description || "No visual layout description specified."}
              </p>

              <div className="pt-2 border-t border-studio-border/60 flex items-center justify-between text-xs text-studio-muted">
                <span className="font-semibold text-studio-text">Lighting Variants:</span>
                <div className="flex items-center gap-2">
                  <span className="flex items-center gap-1 text-amber-400"><Sun className="w-3 h-3" /> Day</span>
                  <span className="flex items-center gap-1 text-indigo-400"><Moon className="w-3 h-3" /> Night</span>
                  <span className="flex items-center gap-1 text-cyan-400"><CloudRain className="w-3 h-3" /> Rain</span>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full bg-studio-card border border-studio-border rounded-2xl p-8 text-center text-studio-muted text-xs">
            No locations registered yet. Add a location to set the scene.
          </div>
        )}
      </div>
    </div>
  );
}
