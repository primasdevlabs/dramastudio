"use client";

import { useState, useRef } from "react";
import { Play, Pause, Maximize2, Sparkles } from "lucide-react";
import { Paper, Group, ActionIcon, Badge, Text } from "@mantine/core";

export function VerticalVideoPlayer({
  src,
  poster,
  shotTitle = "Shot 01 - Close Up",
}: {
  src?: string;
  poster?: string;
  shotTitle?: string;
}) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);

  const togglePlay = () => {
    if (!videoRef.current) return;
    if (isPlaying) {
      videoRef.current.pause();
      setIsPlaying(false);
    } else {
      videoRef.current.play();
      setIsPlaying(true);
    }
  };

  return (
    <Paper
      radius="xl"
      withBorder
      className="relative overflow-hidden bg-black border-studio-border aspect-[9/16] w-full max-w-[280px] mx-auto shadow-2xl flex flex-col justify-between p-4"
    >
      {src ? (
        <video
          ref={videoRef}
          src={src}
          poster={poster}
          className="absolute inset-0 w-full h-full object-cover"
          loop
          playsInline
        />
      ) : (
        <div className="absolute inset-0 bg-gradient-to-b from-slate-900 via-studio-card to-slate-950 flex flex-col items-center justify-center p-6 text-center space-y-3">
          <div className="w-12 h-12 rounded-full bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400">
            <Sparkles className="w-6 h-6" />
          </div>
          <Text size="xs" fw={700} c="cyan.4">
            9:16 Vertical Render
          </Text>
          <Text size="xs" c="dimmed">
            Wan 2.1 T2V Generation Pending
          </Text>
        </div>
      )}

      {/* Overlay Top Bar */}
      <Group justify="space-between" className="relative z-10">
        <Badge color="cyan" variant="filled" size="xs">
          9:16 VERTICAL
        </Badge>
        <Badge color="emerald" variant="dot" size="xs">
          Wan 2.1 Ready
        </Badge>
      </Group>

      {/* Overlay Bottom Bar Controls */}
      <div className="relative z-10 space-y-2 bg-black/60 backdrop-blur-md p-3 rounded-xl border border-white/10">
        <Text size="xs" fw={700} c="white" className="truncate">
          {shotTitle}
        </Text>
        <Group justify="space-between" align="center">
          <ActionIcon
            variant="filled"
            color="cyan"
            radius="xl"
            size="md"
            onClick={togglePlay}
          >
            {isPlaying ? <Pause size={16} /> : <Play size={16} />}
          </ActionIcon>
          <Text size="xs" c="dimmed">
            00:05 / 1080x1920
          </Text>
        </Group>
      </div>
    </Paper>
  );
}
