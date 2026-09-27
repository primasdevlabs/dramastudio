"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Activity, Play, CheckCircle2, XCircle, Pause, RotateCcw, Shield, Cpu } from "lucide-react";
import { api } from "@/lib/api/client";
import { LeadDirectorDecision, ApprovalRequest } from "@/lib/api/types";

export default function ProductionControlTowerPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const { data: decisions } = useQuery({
    queryKey: ["decisions", projectId],
    queryFn: async () => {
      const res = await api.get<{ decisions: LeadDirectorDecision[] }>(`/v1/agents/decisions?project_id=${projectId}`);
      return res.decisions || [];
    },
  });

  const { data: approvals } = useQuery({
    queryKey: ["approvals", projectId],
    queryFn: async () => {
      const res = await api.get<{ approvals: ApprovalRequest[] }>(`/v1/production/approval?project_id=${projectId}`);
      return res.approvals || [];
    },
  });

  const triggerDirectorMutation = useMutation({
    mutationFn: async () => {
      return api.post<LeadDirectorDecision>("/v1/agents/lead-director/step", {
        project_id: projectId,
        episode_id: "ep_001",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["decisions", projectId] });
    },
  });

  const approvalMutation = useMutation({
    mutationFn: async (decision: "APPROVE" | "REJECT" | "REQUEST_REVISION") => {
      return api.post<ApprovalRequest>("/v1/production/approval", {
        project_id: projectId,
        episode_id: "ep_001",
        stage: "Episode",
        target_id: "ep_001",
        decision,
        notes: "Decision submitted via Production Control Tower",
        decided_by: "Lead Producer",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["approvals", projectId] });
    },
  });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between bg-studio-card border border-studio-border p-6 rounded-2xl shadow-xl">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400">
            <Activity className="w-5 h-5 animate-pulse" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Production Control Tower</h1>
            <p className="text-xs text-studio-muted">Lead Director console, active Temporal workflows, and human approval gates</p>
          </div>
        </div>

        <button
          onClick={() => triggerDirectorMutation.mutate()}
          disabled={triggerDirectorMutation.isPending}
          className="inline-flex items-center gap-2 bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white font-bold text-xs px-4 py-2.5 rounded-lg shadow-lg shadow-cyan-500/20 transition-all"
        >
          <Cpu className="w-4 h-4" />
          {triggerDirectorMutation.isPending ? "Executing Step..." : "Trigger Lead Director Loop"}
        </button>
      </div>

      {/* Human Approval Gate Box */}
      <div className="bg-studio-card border border-amber-500/30 rounded-2xl p-6 space-y-4 shadow-xl">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="w-2.5 h-2.5 rounded-full bg-amber-400 animate-ping"></span>
            <h2 className="text-sm font-bold text-amber-400 uppercase tracking-wider">Human Approval Gate</h2>
          </div>
          <span className="text-xs text-studio-muted">Temporal Workflow Paused</span>
        </div>

        <p className="text-xs text-studio-text">
          Episode 01 Storyboard & Wan Video Generations require human approval before advancing to FFmpeg Postproduction Assembly.
        </p>

        <div className="flex items-center gap-3 pt-2">
          <button
            onClick={() => approvalMutation.mutate("APPROVE")}
            disabled={approvalMutation.isPending}
            className="inline-flex items-center gap-1.5 bg-emerald-500 hover:bg-emerald-400 text-black font-bold text-xs px-4 py-2 rounded-lg transition-all"
          >
            <CheckCircle2 className="w-4 h-4" /> Approve Production
          </button>
          <button
            onClick={() => approvalMutation.mutate("REQUEST_REVISION")}
            disabled={approvalMutation.isPending}
            className="inline-flex items-center gap-1.5 bg-amber-500/20 hover:bg-amber-500/30 text-amber-300 font-semibold text-xs px-4 py-2 rounded-lg border border-amber-500/30 transition-all"
          >
            <RotateCcw className="w-4 h-4" /> Request Revision
          </button>
          <button
            onClick={() => approvalMutation.mutate("REJECT")}
            disabled={approvalMutation.isPending}
            className="inline-flex items-center gap-1.5 bg-rose-500/20 hover:bg-rose-500/30 text-rose-300 font-semibold text-xs px-4 py-2 rounded-lg border border-rose-500/30 transition-all"
          >
            <XCircle className="w-4 h-4" /> Reject & Regenerate
          </button>
        </div>
      </div>

      {/* Logs & Decision Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Lead Director Decision Records */}
        <div className="bg-studio-card border border-studio-border rounded-2xl p-5 space-y-4">
          <h3 className="text-sm font-bold text-white flex items-center gap-2">
            <Cpu className="w-4 h-4 text-cyan-400" /> Auditable Lead Director Decisions
          </h3>

          <div className="space-y-3">
            {decisions && decisions.length > 0 ? (
              decisions.map((d) => (
                <div key={d.id} className="bg-studio-panel border border-studio-border p-3.5 rounded-xl space-y-1 text-xs">
                  <div className="flex items-center justify-between">
                    <span className="font-bold text-cyan-400">{d.decision}</span>
                    <span className="text-studio-muted font-mono">{d.decision_maker}</span>
                  </div>
                  <p className="text-studio-text">{d.reason}</p>
                </div>
              ))
            ) : (
              <div className="text-xs text-studio-muted italic p-4 text-center">
                No decisions recorded yet. Click "Trigger Lead Director Loop" above.
              </div>
            )}
          </div>
        </div>

        {/* Approval History */}
        <div className="bg-studio-card border border-studio-border rounded-2xl p-5 space-y-4">
          <h3 className="text-sm font-bold text-white flex items-center gap-2">
            <Shield className="w-4 h-4 text-emerald-400" /> Approval History
          </h3>

          <div className="space-y-3">
            {approvals && approvals.length > 0 ? (
              approvals.map((a) => (
                <div key={a.id} className="bg-studio-panel border border-studio-border p-3.5 rounded-xl space-y-1 text-xs">
                  <div className="flex items-center justify-between">
                    <span className="font-bold text-emerald-400">{a.decision}</span>
                    <span className="text-studio-muted">{a.decided_by}</span>
                  </div>
                  <p className="text-studio-text">{a.notes}</p>
                </div>
              ))
            ) : (
              <div className="text-xs text-studio-muted italic p-4 text-center">
                No human approvals recorded yet.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
