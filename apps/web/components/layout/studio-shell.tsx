"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
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
  FolderPlus,
  Search,
  Activity,
  ChevronRight,
  Menu,
  Sliders,
  Cpu,
} from "lucide-react";
import { Tooltip, Badge, Group, Text, ActionIcon, ThemeIcon, UnstyledButton } from "@mantine/core";
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
    <div className="flex h-screen overflow-hidden bg-studio-bg text-studio-text font-sans">
      {/* Editorial Workstation Sidebar */}
      <aside
        className={`${
          isSidebarOpen ? "w-64" : "w-16"
        } transition-all duration-200 bg-studio-card border-r border-studio-border flex flex-col z-20`}
      >
        {/* Header Logo */}
        <div className="h-14 flex items-center justify-between px-4 border-b border-studio-border">
          <Link href="/projects" className="flex items-center gap-3 overflow-hidden">
            <ThemeIcon color="terracotta" variant="filled" size="md" radius="sm">
              <Clapperboard size={18} className="text-white" />
            </ThemeIcon>
            {isSidebarOpen && (
              <Text fw={800} size="md" className="tracking-wider text-white">
                DRAMA<span className="text-studio-accent">STUDIO</span>
              </Text>
            )}
          </Link>
          <ActionIcon variant="subtle" color="gray" onClick={toggleSidebar} radius="sm" size="sm">
            <Menu size={16} />
          </ActionIcon>
        </div>

        {/* Navigation Items */}
        <div className="flex-1 py-3 overflow-y-auto px-2 space-y-1">
          {projectId && isSidebarOpen && (
            <Text size="xs" fw={700} c="dimmed" tt="uppercase" className="px-3 pb-2 tracking-wider">
              Production Workflow
            </Text>
          )}

          {projectId
            ? NAV_ITEMS.map((item) => {
                const fullPath = `${baseUrl}${item.path}`;
                const isActive =
                  item.path === ""
                    ? pathname === baseUrl
                    : pathname.startsWith(fullPath);
                const Icon = item.icon;

                const content = (
                  <UnstyledButton
                    component={Link}
                    key={item.label}
                    href={fullPath}
                    className={`flex items-center gap-3 w-full px-3 py-2 rounded-md text-xs font-semibold transition-all ${
                      isActive
                        ? "bg-studio-panel text-studio-accent border border-studio-border"
                        : "text-studio-muted hover:text-studio-text hover:bg-studio-panel"
                    }`}
                  >
                    <Icon size={16} className={`shrink-0 ${isActive ? "text-studio-accent" : ""}`} />
                    {isSidebarOpen && <span>{item.label}</span>}
                  </UnstyledButton>
                );

                return isSidebarOpen ? (
                  content
                ) : (
                  <Tooltip label={item.label} position="right" key={item.label} withArrow>
                    {content}
                  </Tooltip>
                );
              })
            : (
              <UnstyledButton
                component={Link}
                href="/projects"
                className="flex items-center gap-3 w-full px-3 py-2 rounded-md text-xs font-semibold bg-studio-panel text-studio-accent border border-studio-border"
              >
                <FolderPlus size={16} className="shrink-0 text-studio-accent" />
                {isSidebarOpen && <span>All Productions</span>}
              </UnstyledButton>
            )}
        </div>

        {/* Footer Settings */}
        <div className="p-3 border-t border-studio-border">
          <UnstyledButton
            component={Link}
            href={projectId ? `/projects/${projectId}/settings` : "/settings"}
            className="flex items-center gap-3 w-full px-3 py-2 rounded-md text-xs font-medium text-studio-muted hover:text-studio-text hover:bg-studio-panel transition-colors"
          >
            <Settings size={16} className="shrink-0" />
            {isSidebarOpen && <span>Studio Settings</span>}
          </UnstyledButton>
        </div>
      </aside>

      {/* Main Workspace Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Editorial Top Header */}
        <header className="h-14 bg-studio-card border-b border-studio-border flex items-center justify-between px-6 z-10">
          <Group gap="md">
            <UnstyledButton
              onClick={toggleCommandPalette}
              className="flex items-center gap-2 bg-studio-panel border border-studio-border px-3 py-1.5 rounded-md text-xs text-studio-muted hover:text-studio-text transition-all"
            >
              <Search size={14} className="text-studio-accent" />
              <span>Command Palette (Ctrl + K)</span>
            </UnstyledButton>

            {projectId && (
              <Group gap="xs">
                <ChevronRight size={14} className="text-studio-muted" />
                <Badge color="terracotta" variant="outline" size="xs">
                  MONITORED MODE
                </Badge>
                <Badge color="green" variant="dot" size="xs">
                  Temporal Active
                </Badge>
              </Group>
            )}
          </Group>

          <Group gap="md">
            <Badge
              leftSection={<Cpu size={12} className="text-studio-accent" />}
              variant="outline"
              color="gray"
              size="sm"
              radius="sm"
            >
              Provider: Wan 2.1 Video
            </Badge>
            <ThemeIcon color="terracotta" variant="filled" size="sm" radius="sm">
              <Text size="xs" fw={700} c="white">
                DS
              </Text>
            </ThemeIcon>
          </Group>
        </header>

        {/* Workspace Content */}
        <main className="flex-1 overflow-y-auto p-6 bg-studio-bg">{children}</main>
      </div>
    </div>
  );
}
