"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Clapperboard,
  BookOpen,
  Users,
  Globe,
  Film,
  Sparkles,
  Layers,
  ShieldCheck,
  Cpu,
  Share2,
  Settings,
  FolderPlus,
  Search,
  Activity,
  ChevronRight,
  Menu,
  Sliders,
} from "lucide-react";
import { useStudioStore } from "@/stores/studio-store";

const NAV_ITEMS = [
  { label: "Overview", path: "", icon: Clapperboard },
  { label: "Series Bible", path: "/bible", icon: BookOpen },
  { label: "Story Workspace", path: "/story", icon: Layers },
  { label: "Characters", path: "/characters", icon: Users },
  { label: "World Studio", path: "/world", icon: Globe },
  { label: "Canon Facts", path: "/canon", icon: ShieldCheck },
  { label: "Storyboard & Shots", path: "/storyboard", icon: Film },
  { label: "Production Tower", path: "/production", icon: Activity },
  { label: "Continuity Center", path: "/continuity", icon: ShieldCheck },
  { label: "Postproduction", path: "/postproduction", icon: Sliders },
  { label: "Publishing", path: "/publishing", icon: Share2 },
];

export function StudioShell({
  children,
  projectId,
}: {
  children: React.ReactNode;
  projectId?: string;
}) {
  const pathname = usePathname();
  const { isSidebarOpen, toggleSidebar, toggleCommandPalette } = useStudioStore();

  const baseUrl = projectId ? `/projects/${projectId}` : "/projects";

  return (
    <div className="flex h-screen overflow-hidden bg-studio-bg text-studio-text">
      {/* Sidebar */}
      <aside
        className={`${
          isSidebarOpen ? "w-64" : "w-16"
        } transition-all duration-200 bg-studio-card border-r border-studio-border flex flex-col z-20`}
      >
        {/* Logo & Toggle */}
        <div className="h-16 flex items-center justify-between px-4 border-b border-studio-border">
          <Link href="/projects" className="flex items-center gap-3 overflow-hidden">
            <div className="w-9 h-9 rounded-lg bg-gradient-to-tr from-cyan-500 to-indigo-600 flex items-center justify-center shrink-0 shadow-lg shadow-cyan-500/20">
              <Clapperboard className="w-5 h-5 text-white" />
            </div>
            {isSidebarOpen && (
              <span className="font-bold text-lg tracking-wider text-white">
                DRAMA<span className="text-cyan-400">STUDIO</span>
              </span>
            )}
          </Link>
          <button
            onClick={toggleSidebar}
            className="p-1.5 rounded-md hover:bg-studio-panel text-studio-muted hover:text-studio-text transition-colors"
          >
            <Menu className="w-4 h-4" />
          </button>
        </div>

        {/* Navigation Items */}
        <div className="flex-1 py-4 overflow-y-auto px-2 space-y-1">
          {projectId && (
            <div className="px-3 pb-2 text-xs font-semibold text-studio-muted uppercase tracking-wider">
              {isSidebarOpen && "Production Workflow"}
            </div>
          )}

          {projectId
            ? NAV_ITEMS.map((item) => {
                const fullPath = `${baseUrl}${item.path}`;
                const isActive =
                  item.path === ""
                    ? pathname === baseUrl
                    : pathname.startsWith(fullPath);
                const Icon = item.icon;

                return (
                  <Link
                    key={item.label}
                    href={fullPath}
                    className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all ${
                      isActive
                        ? "bg-cyan-500/10 text-cyan-400 border border-cyan-500/30"
                        : "text-studio-muted hover:text-studio-text hover:bg-studio-panel"
                    }`}
                  >
                    <Icon className={`w-4 h-4 shrink-0 ${isActive ? "text-cyan-400" : ""}`} />
                    {isSidebarOpen && <span>{item.label}</span>}
                  </Link>
                );
              })
            : (
              <Link
                href="/projects"
                className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium bg-cyan-500/10 text-cyan-400 border border-cyan-500/30"
              >
                <FolderPlus className="w-4 h-4 shrink-0 text-cyan-400" />
                {isSidebarOpen && <span>All Productions</span>}
              </Link>
            )}
        </div>

        {/* Footer Settings */}
        <div className="p-3 border-t border-studio-border">
          <Link
            href={projectId ? `/projects/${projectId}/settings` : "/settings"}
            className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-studio-muted hover:text-studio-text hover:bg-studio-panel transition-colors"
          >
            <Settings className="w-4 h-4 shrink-0" />
            {isSidebarOpen && <span>Studio Settings</span>}
          </Link>
        </div>
      </aside>

      {/* Main Workspace Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Studio Top Header */}
        <header className="h-16 bg-studio-card border-b border-studio-border flex items-center justify-between px-6 z-10">
          <div className="flex items-center gap-4">
            <button
              onClick={toggleCommandPalette}
              className="flex items-center gap-2 bg-studio-panel border border-studio-border px-3 py-1.5 rounded-lg text-xs text-studio-muted hover:text-studio-text hover:border-cyan-500/40 transition-all"
            >
              <Search className="w-3.5 h-3.5 text-cyan-400" />
              <span>Command Palette (Ctrl + K)</span>
            </button>

            {projectId && (
              <div className="flex items-center gap-2 text-xs text-studio-muted">
                <ChevronRight className="w-3.5 h-3.5" />
                <span className="px-2 py-0.5 rounded bg-cyan-500/10 text-cyan-400 font-semibold border border-cyan-500/20">
                  MONITORED MODE
                </span>
                <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 font-medium border border-emerald-500/20 flex items-center gap-1">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
                  Temporal Active
                </span>
              </div>
            )}
          </div>

          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 bg-studio-panel px-3 py-1 rounded-full border border-studio-border text-xs">
              <Cpu className="w-3.5 h-3.5 text-cyan-400" />
              <span className="text-studio-muted">AI Capabilities:</span>
              <span className="font-semibold text-studio-text">Wan 2.1 Video Provider</span>
            </div>
            <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-cyan-500 to-indigo-600 flex items-center justify-center text-white font-bold text-xs shadow-md">
              DS
            </div>
          </div>
        </header>

        {/* Workspace Content */}
        <main className="flex-1 overflow-y-auto p-6 bg-studio-bg">{children}</main>
      </div>
    </div>
  );
}
