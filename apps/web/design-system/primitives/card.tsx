"use client";

import { Card as MantineCard, CardProps } from "@mantine/core";

export function Card({ className = "", children, ...props }: CardProps & { className?: string }) {
  return (
    <MantineCard
      p="lg"
      radius="lg"
      withBorder
      className={`bg-studio-card border-studio-border hover:border-cyan-500/40 transition-all ${className}`}
      {...props}
    >
      {children}
    </MantineCard>
  );
}
