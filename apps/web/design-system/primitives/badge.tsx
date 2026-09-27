"use client";

import { Badge as MantineBadge, BadgeProps, Tooltip as MantineTooltip, TooltipProps } from "@mantine/core";

export function Badge(props: BadgeProps) {
  return <MantineBadge radius="md" size="sm" {...props} />;
}

export function Tooltip(props: TooltipProps) {
  return <MantineTooltip radius="md" withArrow {...props} />;
}
