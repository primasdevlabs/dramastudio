import Link from "next/link";
import { Clapperboard, Clock, DollarSign, ChevronRight, Activity } from "lucide-react";
import { Project } from "@/lib/api/types";

export function ProjectCard({ project }: { project: Project }) {
  return (
    <div className="bg-studio-card border border-studio-border rounded-xl p-5 hover:border-cyan-500/50 transition-all hover:shadow-xl hover:shadow-cyan-500/5 flex flex-col justify-between">
      <div>
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center text-cyan-400">
              <Clapperboard className="w-4 h-4" />
            </div>
            <span className="text-xs font-semibold text-cyan-400 uppercase tracking-wider bg-cyan-500/10 px-2 py-0.5 rounded border border-cyan-500/20">
              {project.genre || "Drama"}
            </span>
          </div>
          <span className="text-xs font-semibold px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 flex items-center gap-1">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
            {project.status || "ACTIVE"}
          </span>
        </div>

        <h3 className="text-lg font-bold text-white mb-1 tracking-tight">{project.name}</h3>
        <p className="text-xs text-studio-muted line-clamp-2 mb-4">
          {project.description || "Serialized AI Drama Production"}
        </p>
      </div>

      <div className="border-t border-studio-border/60 pt-4 flex items-center justify-between mt-2">
        <div className="flex items-center gap-3 text-xs text-studio-muted">
          <span className="flex items-center gap-1">
            <Activity className="w-3.5 h-3.5 text-cyan-400" />
            {project.mode.toUpperCase()}
          </span>
          <span className="flex items-center gap-1">
            <DollarSign className="w-3.5 h-3.5 text-emerald-400" />
            ${project.budget?.current_spent || 0}
          </span>
        </div>

        <Link
          href={`/projects/${project.id}`}
          className="inline-flex items-center gap-1 text-xs font-semibold text-cyan-400 hover:text-cyan-300 transition-colors"
        >
          Open Studio <ChevronRight className="w-3.5 h-3.5" />
        </Link>
      </div>
    </div>
  );
}
