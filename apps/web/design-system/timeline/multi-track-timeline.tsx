"use client";

import { Paper, Group, Stack, Text, Badge, ActionIcon } from "@mantine/core";
import { Film, Music, Mic, Captions, Play, Pause } from "lucide-react";
import { useState } from "react";

export interface TimelineTrack {
  id: string;
  name: string;
  type: "video" | "audio" | "sfx" | "subtitle";
  clips: { id: string; title: string; duration: string; start: string }[];
}

export function MultiTrackTimeline({
  tracks = [
    {
      id: "tr_video",
      name: "Video 9:16 (Wan 2.1)",
      type: "video",
      clips: [{ id: "c1", title: "Shot 001 - Wan T2V", duration: "00:05", start: "00:00" }],
    },
    {
      id: "tr_dialogue",
      name: "Dialogue (ElevenLabs)",
      type: "audio",
      clips: [{ id: "c2", title: "Sarah Voice Line 1", duration: "00:04", start: "00:00" }],
    },
    {
      id: "tr_music",
      name: "Background Score",
      type: "sfx",
      clips: [{ id: "c3", title: "Suspense Atmosphere", duration: "00:05", start: "00:00" }],
    },
  ],
}: {
  tracks?: TimelineTrack[];
}) {
  const [isPlaying, setIsPlaying] = useState(false);

  return (
    <Paper p="md" radius="md" withBorder className="bg-studio-card border-studio-border">
      <Stack gap="md">
        <Group justify="space-between">
          <Group gap="xs">
            <ActionIcon
              variant="filled"
              color="terracotta"
              radius="sm"
              size="md"
              onClick={() => setIsPlaying(!isPlaying)}
            >
              {isPlaying ? <Pause size={16} /> : <Play size={16} />}
            </ActionIcon>
            <Text fw={700} size="sm" c="white">
              Multi-Track FFmpeg Timeline
            </Text>
          </Group>
          <Badge color="terracotta" variant="light" className="font-mono">
            9:16 Vertical (1080x1920)
          </Badge>
        </Group>

        <Stack gap="xs">
          {tracks.map((track) => (
            <div key={track.id} className="bg-studio-panel border border-studio-border rounded-md p-3 flex items-center gap-4">
              <div className="w-48 shrink-0 flex items-center gap-2">
                {track.type === "video" && <Film size={14} className="text-studio-accent" />}
                {track.type === "audio" && <Mic size={14} className="text-amber-400" />}
                {track.type === "sfx" && <Music size={14} className="text-emerald-400" />}
                {track.type === "subtitle" && <Captions size={14} className="text-studio-muted" />}
                <Text size="xs" fw={700} c="white" className="truncate">
                  {track.name}
                </Text>
              </div>

              <div className="flex-1 bg-studio-canvas h-10 rounded border border-studio-border/60 relative overflow-hidden flex items-center px-2">
                {track.clips.map((clip) => (
                  <div
                    key={clip.id}
                    className="bg-studio-accent/20 border border-studio-accent/40 rounded px-3 py-1 text-xs text-studio-text font-mono flex items-center gap-2"
                  >
                    <span>{clip.title}</span>
                    <span className="text-[10px] text-studio-muted font-mono">({clip.duration})</span>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </Stack>
      </Stack>
    </Paper>
  );
}
