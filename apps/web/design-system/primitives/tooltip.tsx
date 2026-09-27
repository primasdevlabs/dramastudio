"use client";

import { Tooltip as MantineTooltip, TooltipProps } from "@mantine/core";

export function Tooltip(props: TooltipProps) {
  return <MantineTooltip radius="md" withArrow {...props} />;
}

export type { TooltipProps };
