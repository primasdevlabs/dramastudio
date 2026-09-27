"use client";

import { Badge as MantineBadge, BadgeProps } from "@mantine/core";

export function Badge(props: BadgeProps) {
  return <MantineBadge radius="md" size="sm" {...props} />;
}
