"use client";

import { useState, useRef } from "react";
import { Play, Pause, Film } from "lucide-react";
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
      radius="md"
      withBorder
      className="relative overflow-hidden bg-black border-studio-border aspect-[9/16] w-full max-w-[280px] mx-auto shadow-2xl flex flex-col justify-between p-3"
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
        <div className="absolute inset-0 bg-studio-surface flex flex-col items-center justify-center p-6 text-center space-y-3">
          <div className="w-10 h-10 rounded-sm bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
            <Film className="w-5 h-5" />
          </div>
          <Text size="xs" fw={700} c="terracotta.6">
            9:16 Vertical Monitor
          </Text>
          <Text size="xs" c="dimmed">
            Generation: Wan 2.1 Video (T2V)
          </Text>
        </div>
      )}

      {/* Overlay Top Bar */}
      <Group justify="space-between" className="relative z-10">
        <Badge color="terracotta" variant="filled" size="xs">
          9:16 VERTICAL
        </Badge>
        <Badge color="green" variant="outline" size="xs">
          Wan 2.1 Ready
        </Badge>
      </Group>

      {/* Overlay Bottom Technical Metadata Bar */}
      <div className="relative z-10 space-y-2 bg-studio-surface/90 p-2.5 rounded-sm border border-studio-border">
        <Text size="xs" fw={700} c="white" className="truncate">
          {shotTitle}
        </Text>
        <Group justify="space-between" align="center">
          <ActionIcon
            variant="filled"
            color="terracotta"
            radius="sm"
            size="sm"
            onClick={togglePlay}
          >
            {isPlaying ? <Pause size={14} /> : <Play size={14} />}
          </ActionIcon>
          <Text size="xs" c="dimmed" className="font-mono">
            00:04.82 / 1080x1920
          </Text>
        </Group>
      </div>
    </Paper>
  );
}
