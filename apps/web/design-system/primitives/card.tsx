"use client";

import { Card as MantineCard, CardProps } from "@mantine/core";

export function Card({ className = "", children, ...props }: CardProps & { className?: string }) {
  return (
    <MantineCard
      p="lg"
      radius="md"
      withBorder
      className={`bg-studio-card border-studio-border hover:border-studio-accent/40 transition-all ${className}`}
      {...props}
    >
      {children}
    </MantineCard>
  );
}
