"use client";

import Link from "next/link";
import { usePathname, useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
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
  ChevronLeft,
  ChevronDown,
  Calendar,
  PanelLeftClose,
  PanelLeftOpen,
  Sliders,
  Cpu,
  Bell,
  User,
  LogOut,
  HelpCircle,
  Command,
  FolderKanban,
} from "lucide-react";
import {
  Tooltip,
  Badge,
  Group,
  Text,
  ActionIcon,
  ThemeIcon,
  UnstyledButton,
  Popover,
  Stack,
  Divider,
} from "@mantine/core";
import { useStudioStore } from "@/stores/studio-store";
import { CommandPalette } from "./command-palette";

// Global Studio Navigation (active when no project is selected)
const GLOBAL_NAV_ITEMS = [
  { label: "All Productions", path: "/projects", icon: FolderKanban },
  { label: "New Production", path: "/projects/new", icon: FolderPlus },
  { label: "Global Settings", path: "/settings", icon: Settings },
];

// Project Workspace Navigation (active within a specific project)
const PROJECT_NAV_ITEMS = [
  { label: "Overview", path: "", icon: Clapperboard },
  { label: "Series Bible", path: "/bible", icon: BookOpen },
  { label: "Story Workspace", path: "/story", icon: Layers },
  { label: "Seasons & Episodes", path: "/seasons", icon: Calendar },
  { label: "Characters", path: "/characters", icon: Users },
  { label: "World Studio", path: "/world", icon: Globe },
  { label: "Canon Facts", path: "/canon", icon: ShieldCheck },
  { label: "Storyboard & Shots", path: "/storyboard", icon: Film },
  { label: "Production Tower", path: "/production", icon: Activity },
  { label: "Continuity Center", path: "/continuity", icon: ShieldCheck },
  { label: "Postproduction", path: "/postproduction", icon: Sliders },
  { label: "Publishing", path: "/publishing", icon: Share2 },
  { label: "Project Settings", path: "/settings", icon: Settings },
];

function TopbarSearch({ onOpen }: { onOpen: () => void }) {
  return (
    <UnstyledButton
      onClick={onOpen}
      className="flex items-center gap-2 bg-studio-panel border border-studio-border px-2.5 py-1 rounded-md text-[11px] text-studio-muted hover:text-studio-text hover:border-studio-accent/40 transition-all w-52"
      aria-label="Open command palette"
    >
      <Search size={13} className="text-studio-accent shrink-0" />
      <span className="flex-1 text-left text-[11px]">Search or command...</span>
      <kbd className="hidden sm:inline-flex items-center gap-0.5 px-1 py-0.5 rounded bg-studio-bg border border-studio-border text-[9px] text-studio-muted font-normal">
        <Command size={9} /> K
      </kbd>
    </UnstyledButton>
  );
}

function NotificationDropdown() {
  const [opened, setOpened] = useState(false);

  return (
    <Popover
      opened={opened}
      onChange={setOpened}
      position="bottom-end"
      width={280}
      shadow="md"
      withArrow
    >
      <Popover.Target>
        <ActionIcon
          variant="subtle"
          color="gray"
          radius="sm"
          size="sm"
          onClick={() => setOpened(!opened)}
          aria-label="Notifications"
          className="relative"
        >
          <Bell size={14} />
          <span className="absolute top-0.5 right-0.5 w-1.5 h-1.5 rounded-full bg-studio-accent" />
        </ActionIcon>
      </Popover.Target>
      <Popover.Dropdown className="bg-studio-card border border-studio-border p-0">
        <div className="px-3 py-2 border-b border-studio-border">
          <Text fw={600} fz={11} c="white">
            Notifications
          </Text>
        </div>
        <Stack gap={0} className="max-h-72 overflow-y-auto">
          {[
            {
              title: "Generation complete",
              desc: "Shot 07 video generated",
              time: "2m ago",
              read: false,
            },
            {
              title: "Continuity alert",
              desc: "Wardrobe mismatch in Scene 04",
              time: "12m ago",
              read: false,
            },
            {
              title: "Approval gate",
              desc: "Episode 12 storyboard ready for review",
              time: "1h ago",
              read: true,
            },
          ].map((n, i) => (
            <UnstyledButton
              key={i}
              className={`px-3 py-2 hover:bg-studio-panel transition-colors border-b border-studio-border/50 ${
                !n.read ? "bg-studio-accent/5" : ""
              }`}
            >
              <div className="flex items-start gap-2">
                {!n.read && (
                  <span className="w-1.5 h-1.5 rounded-full bg-studio-accent mt-1 shrink-0" />
                )}
                <div>
                  <Text fw={600} fz={11} c="white">
                    {n.title}
                  </Text>
                  <Text fz={10} c="dimmed">
                    {n.desc}
                  </Text>
                  <Text fz={9} c="dimmed" mt={1}>
                    {n.time}
                  </Text>
                </div>
              </div>
            </UnstyledButton>
          ))}
        </Stack>
        <div className="px-3 py-1.5 border-t border-studio-border text-center">
          <UnstyledButton className="text-[10px] text-studio-accent hover:underline">
            View all notifications
          </UnstyledButton>
        </div>
      </Popover.Dropdown>
    </Popover>
  );
}

function ProfileDropdown({ projectId }: { projectId?: string }) {
  const [opened, setOpened] = useState(false);
  const router = useRouter();

  return (
    <Popover
      opened={opened}
      onChange={setOpened}
      position="bottom-end"
      width={200}
      shadow="md"
      withArrow
    >
      <Popover.Target>
        <UnstyledButton
          onClick={() => setOpened(!opened)}
          className="flex items-center gap-1.5 px-1.5 py-1 rounded-md hover:bg-studio-panel transition-colors"
          aria-label="User menu"
        >
          <ThemeIcon color="terracotta" variant="filled" size={18} radius="xl">
            <Text fz={9} fw={700} c="white">
              DS
            </Text>
          </ThemeIcon>
          <ChevronDown size={11} className="text-studio-muted" />
        </UnstyledButton>
      </Popover.Target>
      <Popover.Dropdown className="bg-studio-card border border-studio-border p-0 rounded-md shadow-xl">
        <div className="px-3 py-2 border-b border-studio-border">
          <Text fw={600} fz={11} c="white">
            Studio Operator
          </Text>
          <Text fz={10} c="dimmed" mt={0.5}>
            operator@dramastudio.ai
          </Text>
        </div>
        <div className="p-1 space-y-0.5">
          <UnstyledButton
            component={Link}
            href="/settings"
            onClick={() => setOpened(false)}
            className="flex items-center gap-2 px-2.5 py-1.5 rounded hover:bg-studio-panel transition-colors w-full"
          >
            <Globe size={12} className="text-studio-accent shrink-0" />
            <Text fz={10} fw={500} c="dimmed" className="hover:text-white">
              Global Studio Settings
            </Text>
          </UnstyledButton>

          {projectId && (
            <UnstyledButton
              component={Link}
              href={`/projects/${projectId}/settings`}
              onClick={() => setOpened(false)}
              className="flex items-center gap-2 px-2.5 py-1.5 rounded hover:bg-studio-panel transition-colors w-full"
            >
              <Settings size={12} className="text-studio-muted shrink-0" />
              <Text fz={10} fw={500} c="dimmed" className="hover:text-white">
                Project Settings
              </Text>
            </UnstyledButton>
          )}

          <UnstyledButton
            onClick={() => setOpened(false)}
            className="flex items-center gap-2 px-2.5 py-1.5 rounded hover:bg-studio-panel transition-colors w-full"
          >
            <HelpCircle size={12} className="text-studio-muted shrink-0" />
            <Text fz={10} fw={500} c="dimmed" className="hover:text-white">
              Documentation
            </Text>
          </UnstyledButton>
        </div>
        <Divider color="dark.5" />
        <div className="p-1">
          <UnstyledButton
            onClick={() => {
              setOpened(false);
              router.push("/projects");
            }}
            className="flex items-center gap-2 px-2.5 py-1.5 rounded hover:bg-studio-panel transition-colors w-full"
          >
            <LogOut size={12} className="text-studio-red shrink-0" />
            <Text fz={10} fw={500} c="red">
              Sign Out
            </Text>
          </UnstyledButton>
        </div>
      </Popover.Dropdown>
    </Popover>
  );
}

function Breadcrumbs({ projectId }: { projectId?: string }) {
  const pathname = usePathname();

  if (!projectId) return null;

  const segments = pathname
    .replace(`/projects/${projectId}`, "")
    .split("/")
    .filter(Boolean);

  if (segments.length === 0) return null;

  const crumbs = [
    { label: "Projects", href: "/projects" },
    ...segments.map((seg, i) => ({
      label: seg
        .replace(/-/g, " ")
        .replace(/\b\w/g, (c) => c.toUpperCase()),
      href: `/projects/${projectId}/${segments.slice(0, i + 1).join("/")}`,
    })),
  ];

  return (
    <nav aria-label="Breadcrumb" className="flex items-center gap-1">
      {crumbs.map((crumb, i) => (
        <span key={crumb.href} className="flex items-center gap-1">
          {i > 0 && (
            <ChevronRight size={11} className="text-studio-muted" />
          )}
          {i === crumbs.length - 1 ? (
            <Text fz={11} fw={500} c="white">
              {crumb.label}
            </Text>
          ) : (
            <Link href={crumb.href}>
              <Text fz={11} c="dimmed" className="hover:text-studio-accent transition-colors">
                {crumb.label}
              </Text>
            </Link>
          )}
        </span>
      ))}
    </nav>
  );
}

export function StudioShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const params = useParams<{ projectId?: string }>();
  const projectId = params?.projectId;

  const {
    isSidebarOpen,
    toggleSidebar,
    isCommandPaletteOpen,
    setCommandPaletteOpen,
  } = useStudioStore();

  // Select appropriate nav items based on context
  const navItems = projectId ? PROJECT_NAV_ITEMS : GLOBAL_NAV_ITEMS;
  const projectBaseUrl = projectId ? `/projects/${projectId}` : "";

  // Global keyboard shortcut: Ctrl+K / Cmd+K
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setCommandPaletteOpen(!isCommandPaletteOpen);
      }
      if (e.key === "Escape" && isCommandPaletteOpen) {
        setCommandPaletteOpen(false);
      }
    }
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isCommandPaletteOpen, setCommandPaletteOpen]);

  return (
    <div className="flex h-screen bg-studio-bg text-studio-text overflow-hidden">
      {/* Left Navigation Sidebar */}
      <aside
        className={`${
          isSidebarOpen ? "w-60" : "w-16"
        } bg-studio-card border-r border-studio-border flex flex-col transition-all duration-200 shrink-0 z-20`}
      >
        {/* Sidebar Header */}
        <div className="h-14 flex items-center justify-between px-4 border-b border-studio-border">
          {isSidebarOpen ? (
            <Link
              href="/projects"
              className="flex items-center gap-2 text-white font-bold text-sm tracking-tight"
            >
              <div className="w-7 h-7 rounded bg-studio-accent flex items-center justify-center text-white font-black text-xs shadow-sm">
                DS
              </div>
              <span>DramaStudio</span>
            </Link>
          ) : (
            <Link href="/projects" className="mx-auto">
              <div className="w-7 h-7 rounded bg-studio-accent flex items-center justify-center text-white font-black text-xs shadow-sm">
                DS
              </div>
            </Link>
          )}

          <ActionIcon
            variant="subtle"
            color="gray"
            size="sm"
            onClick={toggleSidebar}
            aria-label="Toggle sidebar"
            className="text-studio-muted hover:text-studio-text"
          >
            {isSidebarOpen ? (
              <PanelLeftClose size={16} />
            ) : (
              <PanelLeftOpen size={16} />
            )}
          </ActionIcon>
        </div>

        {/* Project Selector / Return Indicator */}
        {projectId ? (
          <div className="px-3 py-2.5 border-b border-studio-border bg-studio-panel/50 space-y-2">
            <Link
              href="/projects"
              className="flex items-center gap-1.5 text-studio-accent hover:underline text-[10px] font-medium"
            >
              <ChevronLeft size={12} />
              <span>Back to All Productions</span>
            </Link>

            {isSidebarOpen && (
              <div className="flex items-center justify-between">
                <Text fw={700} fz={11} c="white" className="truncate">
                  Midnight Call
                </Text>
                <Badge color="terracotta" variant="light" size="xs">
                  S01
                </Badge>
              </div>
            )}
          </div>
        ) : (
          isSidebarOpen && (
            <div className="px-3.5 py-2 border-b border-studio-border">
              <Text fz={10} c="dimmed" fw={600} tt="uppercase" className="tracking-wider">
                GLOBAL STUDIO
              </Text>
            </div>
          )
        )}

        {/* Navigation Items */}
        <nav className="flex-1 overflow-y-auto py-3 px-2 space-y-0.5">
          {navItems.map((item) => {
            const fullPath = projectId
              ? `${projectBaseUrl}${item.path}`
              : item.path;

            const isActive = projectId
              ? item.path === ""
                ? pathname === projectBaseUrl
                : pathname.startsWith(fullPath)
              : pathname === item.path || (item.path !== "/projects" && pathname.startsWith(item.path));

            const Icon = item.icon;

            const linkContent = (
              <Link
                key={item.label}
                href={fullPath}
                className={`flex items-center gap-3 px-3 py-2 rounded-md text-xs font-medium transition-all ${
                  isActive
                    ? "bg-studio-accent text-white shadow-sm font-semibold"
                    : "text-studio-muted hover:text-studio-text hover:bg-studio-panel"
                } ${!isSidebarOpen ? "justify-center px-0" : ""}`}
              >
                <Icon
                  size={16}
                  className={
                    isActive ? "text-white" : "text-studio-muted shrink-0"
                  }
                />
                {isSidebarOpen && <span className="truncate">{item.label}</span>}
              </Link>
            );

            if (!isSidebarOpen) {
              return (
                <Tooltip
                  key={item.label}
                  label={item.label}
                  position="right"
                  withArrow
                >
                  <div>{linkContent}</div>
                </Tooltip>
              );
            }

            return linkContent;
          })}
        </nav>

        {/* Sidebar Footer */}
        <div className="p-3 border-t border-studio-border">
          {isSidebarOpen ? (
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                <Text fz={10} c="dimmed" fw={600} tt="uppercase">
                  Engine Ready
                </Text>
              </div>
              <Badge color="gray" variant="outline" size="xs">
                v2.4
              </Badge>
            </div>
          ) : (
            <div className="w-2 h-2 rounded-full bg-emerald-500 mx-auto" />
          )}
        </div>
      </aside>

      {/* Main Area */}
      <div className="flex-1 flex flex-col overflow-hidden min-w-0">
        {/* Top Header */}
        <header className="h-12 bg-studio-card border-b border-studio-border flex items-center justify-between px-3.5 z-10 shrink-0">
          <Group gap="sm">
            <TopbarSearch onOpen={() => setCommandPaletteOpen(true)} />
            <Divider orientation="vertical" color="dark.5" className="h-4" />
            <Breadcrumbs projectId={projectId} />

            {projectId && (
              <Group gap="xs">
                <Badge color="terracotta" variant="outline" size="xs">
                  MONITORED
                </Badge>
                <Badge color="green" variant="dot" size="xs">
                  Temporal Active
                </Badge>
              </Group>
            )}
          </Group>

          <Group gap="xs">
            <NotificationDropdown />
            <Divider orientation="vertical" color="dark.5" className="h-4" />
            <ProfileDropdown projectId={projectId} />
          </Group>
        </header>

        {/* Workspace Content */}
        <main className="flex-1 overflow-y-auto p-6 bg-studio-bg">
          {children}
        </main>
      </div>

      {/* Command Palette */}
      <CommandPalette
        opened={isCommandPaletteOpen}
        onClose={() => setCommandPaletteOpen(false)}
        projectId={projectId}
      />
    </div>
  );
}
