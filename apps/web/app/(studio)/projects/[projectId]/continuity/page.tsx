"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ShieldCheck, AlertOctagon, AlertTriangle, Info, CheckCircle2 } from "lucide-react";
import { useParams } from "next/navigation";
import { api } from "@/lib/api/client";
import { ContinuityIssue } from "@/lib/api/types";

export default function ContinuityCenterPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const queryClient = useQueryClient();

  const { data: issues, isLoading } = useQuery({
    queryKey: ["continuity-issues", projectId],
    queryFn: async () => {
      const res = await api.get<{ items: ContinuityIssue[] }>(`/v1/projects/${projectId}/continuity/issues`);
      return res.items || [];
    },
  });

  const resolveMutation = useMutation({
    mutationFn: async (issueId: string) => {
      return api.post(`/v1/projects/${projectId}/continuity/issues/${issueId}/resolve`, {
        resolution: "Corrected shot brief to match canon wardrobe specification",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["continuity-issues", projectId] });
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-md">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded bg-rose-500/10 border border-rose-500/30 flex items-center justify-center text-rose-400">
            <ShieldCheck className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Continuity Center</h1>
            <p className="text-xs text-studio-muted">Automated story, character, visual, timeline, and wardrobe continuity verification</p>
          </div>
        </div>
      </div>

      {/* Issues List */}
      <div className="space-y-4">
        {isLoading ? (
          [1, 2].map((i) => <div key={i} className="h-32 bg-studio-card rounded animate-pulse" />)
        ) : issues && issues.length > 0 ? (
          issues.map((issue) => (
            <div
              key={issue.id}
              className={`bg-studio-card border rounded-md p-5 space-y-3 transition-all ${
                issue.status === "resolved" || issue.status === "wontfix"
                  ? "border-emerald-500/30 bg-emerald-500/5"
                  : issue.severity === "BLOCKING"
                  ? "border-rose-500/50 bg-rose-500/5"
                  : "border-amber-500/40 bg-amber-500/5"
              }`}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  {issue.severity === "BLOCKING" ? (
                    <AlertOctagon className="w-4 h-4 text-rose-400" />
                  ) : (
                    <AlertTriangle className="w-4 h-4 text-amber-400" />
                  )}
                  <span className="text-sm font-bold text-white">{issue.category} Violation — {issue.entity}</span>
                </div>
                <div className="flex items-center gap-2">
                  <span
                    className={`px-2.5 py-0.5 rounded text-xs font-bold border ${
                      issue.severity === "BLOCKING"
                        ? "bg-rose-500/10 text-rose-400 border-rose-500/20"
                        : "bg-amber-500/10 text-amber-400 border-amber-500/20"
                    }`}
                  >
                    {issue.severity}
                  </span>
                  {issue.status === "resolved" || issue.status === "wontfix" ? (
                    <span className="px-2.5 py-0.5 rounded bg-emerald-500/10 text-emerald-400 text-xs font-bold border border-emerald-500/20 flex items-center gap-1">
                      <CheckCircle2 className="w-3.5 h-3.5" /> Resolved
                    </span>
                  ) : (
                    <button
                      onClick={() => resolveMutation.mutate(issue.id)}
                      disabled={resolveMutation.isPending}
                      className="bg-emerald-500 hover:bg-emerald-400 text-black font-bold text-xs px-3 py-1 rounded transition-all"
                    >
                      Resolve Issue
                    </button>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs bg-studio-panel p-3.5 rounded border border-studio-border/60">
                <div>
                  <span className="text-studio-muted font-semibold">Expected Canonical State:</span>
                  <p className="text-emerald-400 font-semibold mt-0.5">{issue.expected_state}</p>
                </div>
                <div>
                  <span className="text-studio-muted font-semibold">Actual Generated State:</span>
                  <p className="text-rose-400 font-semibold mt-0.5">{issue.actual_state}</p>
                </div>
              </div>

              {issue.cause && (
                <p className="text-xs text-studio-muted">
                  <span className="font-semibold text-studio-text">Cause & Evidence:</span> {issue.cause} — {issue.evidence}
                </p>
              )}
            </div>
          ))
        ) : (
          <div className="bg-studio-card border border-studio-border rounded-md p-8 text-center text-studio-muted text-xs">
            No continuity issues flagged. All story, wardrobe, and visual assets match canonical truth.
          </div>
        )}
      </div>
    </div>
  );
}
