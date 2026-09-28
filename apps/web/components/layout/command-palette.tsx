"use client";

import { useState, useEffect, useRef, useMemo } from "react";
import { useRouter } from "next/navigation";
import {
  Clapperboard,
  BookOpen,
  Users,
  Globe,
  Film,
  Layers,
  ShieldCheck,
  Share2,
  Settings,
  Activity,
  Calendar,
  Sliders,
  Search,
  ArrowRight,
  Pause,
  Play,
  FileText,
  Image,
  Bot,
} from "lucide-react";
import { Modal, Text, TextInput, Group, Badge, Stack, UnstyledButton, ThemeIcon } from "@mantine/core";

interface CommandItem {
  id: string;
  label: string;
  description: string;
  icon: React.ElementType;
  category: "navigation" | "production" | "search";
  action: () => void;
  keywords: string[];
}

interface CommandPaletteProps {
  opened: boolean;
  onClose: () => void;
  projectId?: string;
}

export function CommandPalette({ opened, onClose, projectId }: CommandPaletteProps) {
  const [query, setQuery] = useState("");
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  const base = projectId ? `/projects/${projectId}` : "/projects";

  const commands = useMemo<CommandItem[]>(() => {
    const nav: CommandItem[] = [
      {
        id: "projects",
        label: "All Projects",
        description: "View all productions",
        icon: Clapperboard,
        category: "navigation",
        action: () => router.push("/projects"),
        keywords: ["projects", "productions", "home"],
      },
      {
        id: "global-settings",
        label: "Global Studio Settings",
        description: "Studio-wide AI model defaults & master API keys",
        icon: Globe,
        category: "navigation",
        action: () => router.push("/settings"),
        keywords: ["global", "settings", "defaults", "api", "keys", "preferences"],
      },
    ];

    if (projectId) {
      nav.push(
        {
          id: "overview",
          label: "Project Overview",
          description: "Dashboard for current project",
          icon: Clapperboard,
          category: "navigation",
          action: () => router.push(base),
          keywords: ["overview", "dashboard"],
        },
        {
          id: "bible",
          label: "Series Bible",
          description: "Premise, tone, world & narrative rules",
          icon: BookOpen,
          category: "navigation",
          action: () => router.push(`${base}/bible`),
          keywords: ["bible", "premise", "narrative"],
        },
        {
          id: "story",
          label: "Story Workspace",
          description: "Story graph, timeline, and arcs",
          icon: Layers,
          category: "navigation",
          action: () => router.push(`${base}/story`),
          keywords: ["story", "graph", "timeline", "arcs"],
        },
        {
          id: "seasons",
          label: "Seasons & Episodes",
          description: "Season arcs and episode planning",
          icon: Calendar,
          category: "navigation",
          action: () => router.push(`${base}/seasons`),
          keywords: ["seasons", "episodes", "planning"],
        },
        {
          id: "characters",
          label: "Character Studio",
          description: "Identity, wardrobe, voice profiles",
          icon: Users,
          category: "navigation",
          action: () => router.push(`${base}/characters`),
          keywords: ["characters", "identity", "wardrobe"],
        },
        {
          id: "world",
          label: "World Studio",
          description: "Locations, environments & props",
          icon: Globe,
          category: "navigation",
          action: () => router.push(`${base}/world`),
          keywords: ["world", "locations", "environments"],
        },
        {
          id: "canon",
          label: "Canon Facts",
          description: "Canonical story facts and rules",
          icon: ShieldCheck,
          category: "navigation",
          action: () => router.push(`${base}/canon`),
          keywords: ["canon", "facts", "rules"],
        },
        {
          id: "storyboard",
          label: "Storyboard & Shots",
          description: "Camera angles, shots & visual panels",
          icon: Film,
          category: "navigation",
          action: () => router.push(`${base}/storyboard`),
          keywords: ["storyboard", "shots", "camera"],
        },
        {
          id: "production",
          label: "Production Console",
          description: "Production activity & stage status",
          icon: Activity,
          category: "navigation",
          action: () => router.push(`${base}/production`),
          keywords: ["production", "console", "status", "roles"],
        },
        {
          id: "continuity",
          label: "Continuity Center",
          description: "Wardrobe & story timeline checks",
          icon: ShieldCheck,
          category: "navigation",
          action: () => router.push(`${base}/continuity`),
          keywords: ["continuity", "wardrobe", "timeline"],
        },
        {
          id: "postproduction",
          label: "Postproduction",
          description: "Assembly timeline, renderer",
          icon: Sliders,
          category: "navigation",
          action: () => router.push(`${base}/postproduction`),
          keywords: ["postproduction", "assembly", "render"],
        },
        {
          id: "publishing",
          label: "Publishing",
          description: "Distribution channels and scheduling",
          icon: Share2,
          category: "navigation",
          action: () => router.push(`${base}/publishing`),
          keywords: ["publishing", "distribute", "tiktok", "youtube"],
        },
        {
          id: "settings",
          label: "Studio Settings",
          description: "Capabilities, policies, integrations",
          icon: Settings,
          category: "navigation",
          action: () => router.push(`${base}/settings`),
          keywords: ["settings", "configuration", "capabilities"],
        },
      );
    }

    if (projectId) {
      nav.push(
        {
          id: "pause",
          label: "Pause Production",
          description: "Halt active production processes",
          icon: Pause,
          category: "production",
          action: () => { onClose(); },
          keywords: ["pause", "stop", "halt"],
        },
        {
          id: "resume",
          label: "Resume Production",
          description: "Resume active production processes",
          icon: Play,
          category: "production",
          action: () => { onClose(); },
          keywords: ["resume", "start", "continue"],
        },
        {
          id: "gen-script",
          label: "Write Script",
          description: "Write episode screenplay from series bible",
          icon: FileText,
          category: "production",
          action: () => { onClose(); },
          keywords: ["write", "script", "screenplay", "dialogue"],
        },
        {
          id: "gen-storyboard",
          label: "Create Storyboard",
          description: "Create storyboard panels & shot list",
          icon: Image,
          category: "production",
          action: () => { onClose(); },
          keywords: ["create", "storyboard", "shots", "framing"],
        },
        {
          id: "agent-activity",
          label: "Production Console & Activity",
          description: "Monitor production roles & automated tasks",
          icon: Bot,
          category: "production",
          action: () => router.push(`${base}/production`),
          keywords: ["production", "activity", "console", "roles"],
        },
      );
    }

    return nav;
  }, [projectId, base, router, onClose]);

  const filtered = useMemo(() => {
    if (!query.trim()) return commands;
    const q = query.toLowerCase();
    return commands.filter(
      (cmd) =>
        cmd.label.toLowerCase().includes(q) ||
        cmd.description.toLowerCase().includes(q) ||
        cmd.keywords.some((kw) => kw.includes(q))
    );
  }, [query, commands]);

  // Reset on open
  useEffect(() => {
    if (opened) {
      setQuery("");
      setSelectedIndex(0);
      setTimeout(() => inputRef.current?.focus(), 50);
    }
  }, [opened]);

  // Keyboard nav
  useEffect(() => {
    if (!opened) return;

    function handleKey(e: KeyboardEvent) {
      if (e.key === "ArrowDown") {
        e.preventDefault();
        setSelectedIndex((i) => Math.min(i + 1, filtered.length - 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setSelectedIndex((i) => Math.max(i - 1, 0));
      } else if (e.key === "Enter" && filtered[selectedIndex]) {
        e.preventDefault();
        filtered[selectedIndex].action();
        onClose();
      }
    }

    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [opened, filtered, selectedIndex, onClose]);

  // Clamp selected index when filtered results change
  useEffect(() => {
    setSelectedIndex((i) => Math.min(i, Math.max(filtered.length - 1, 0)));
  }, [filtered.length]);

  const grouped = useMemo(() => {
    const groups: Record<string, CommandItem[]> = {};
    for (const cmd of filtered) {
      if (!groups[cmd.category]) groups[cmd.category] = [];
      groups[cmd.category].push(cmd);
    }
    return groups;
  }, [filtered]);

  const categoryLabels: Record<string, string> = {
    navigation: "Navigate",
    production: "Production Commands",
    search: "Search",
  };

  let flatIndex = -1;

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      withCloseButton={false}
      centered
      size="lg"
      padding={0}
      radius="md"
      overlayProps={{ backgroundOpacity: 0.6, blur: 4 }}
      classNames={{
        content: "bg-studio-card border border-studio-border overflow-hidden shadow-2xl",
      }}
    >
      {/* Search Input Bar */}
      <div className="flex items-center gap-3 border-b border-studio-border px-4 py-1 bg-studio-panel/50">
        <Search size={16} className="text-studio-accent shrink-0" />
        <input
          ref={inputRef}
          value={query}
          onChange={(e) => setQuery(e.currentTarget.value)}
          placeholder="Search commands, pages, or actions..."
          className="w-full bg-transparent text-sm text-studio-text placeholder:text-studio-muted py-3 outline-none border-none font-medium"
        />
        {query && (
          <button
            onClick={() => setQuery("")}
            className="text-xs text-studio-muted hover:text-studio-text px-2 py-1 rounded bg-studio-bg"
          >
            Clear
          </button>
        )}
      </div>

      {/* Results List */}
      <div className="max-h-80 overflow-y-auto py-2 px-2">
        {filtered.length === 0 ? (
          <div className="px-6 py-10 text-center">
            <Text size="xs" c="dimmed">No results found for &ldquo;{query}&rdquo;</Text>
          </div>
        ) : (
          Object.entries(grouped).map(([category, items]) => (
            <div key={category} className="mb-2">
              <div className="px-3 py-1.5 text-[10px] font-bold text-studio-muted/80 uppercase tracking-widest">
                {categoryLabels[category] ?? category}
              </div>
              <div className="space-y-0.5">
                {items.map((cmd) => {
                  flatIndex++;
                  const isSelected = flatIndex === selectedIndex;
                  const Icon = cmd.icon;
                  const idx = flatIndex;

                  return (
                    <UnstyledButton
                      key={cmd.id}
                      onClick={() => {
                        cmd.action();
                        onClose();
                      }}
                      onMouseEnter={() => setSelectedIndex(idx)}
                      style={{
                        paddingLeft: "1rem",
                        paddingRight: "1rem",
                        paddingTop: "0.625rem",
                        paddingBottom: "0.625rem",
                      }}
                      className={`flex items-center gap-3.5 w-full rounded-md text-xs transition-all ${
                        isSelected
                          ? "bg-studio-panel text-studio-text border border-studio-border/60 shadow-sm"
                          : "text-studio-muted hover:bg-studio-panel/50 hover:text-studio-text"
                      }`}
                    >
                      <Icon
                        size={16}
                        className={`shrink-0 ${
                          isSelected ? "text-studio-accent" : "text-studio-muted"
                        }`}
                      />
                      <div className="flex-1 min-w-0 text-left">
                        <Text
                          size="xs"
                          fw={600}
                          c={isSelected ? "white" : undefined}
                          truncate="end"
                        >
                          {cmd.label}
                        </Text>
                        <Text size="xs" c="dimmed" truncate="end" className="text-[11px]">
                          {cmd.description}
                        </Text>
                      </div>
                      {isSelected && (
                        <ArrowRight size={13} className="text-studio-accent shrink-0 ml-2" />
                      )}
                    </UnstyledButton>
                  );
                })}
              </div>
            </div>
          ))
        )}
      </div>

      {/* Footer Instructions */}
      <div className="border-t border-studio-border px-4 py-2.5 bg-studio-panel/30 flex items-center justify-between text-xs">
        <Group gap="md">
          <Group gap={4}>
            <kbd className="px-1.5 py-0.5 rounded bg-studio-panel border border-studio-border text-[10px] font-sans text-studio-muted">↑↓</kbd>
            <Text size="xs" c="dimmed">navigate</Text>
          </Group>
          <Group gap={4}>
            <kbd className="px-1.5 py-0.5 rounded bg-studio-panel border border-studio-border text-[10px] font-sans text-studio-muted">↵</kbd>
            <Text size="xs" c="dimmed">select</Text>
          </Group>
          <Group gap={4}>
            <kbd className="px-1.5 py-0.5 rounded bg-studio-panel border border-studio-border text-[10px] font-sans text-studio-muted">esc</kbd>
            <Text size="xs" c="dimmed">close</Text>
          </Group>
        </Group>
      </div>
    </Modal>
  );
}
