"use client";

import { Button as MantineButton, ButtonProps as MantineButtonProps } from "@mantine/core";
import { forwardRef } from "react";

export interface StudioButtonProps extends MantineButtonProps {
  glow?: "cyan" | "amber" | "emerald" | "none";
}

export const Button = forwardRef<HTMLButtonElement, StudioButtonProps>(
  ({ glow = "none", className = "", children, ...props }, ref) => {
    let glowClass = "";
    if (glow === "cyan") glowClass = "shadow-lg shadow-cyan-500/20";
    if (glow === "amber") glowClass = "shadow-lg shadow-amber-500/20";
    if (glow === "emerald") glowClass = "shadow-lg shadow-emerald-500/20";

    return (
      <MantineButton
        ref={ref}
        radius="md"
        className={`${glowClass} ${className}`}
        {...props}
      >
        {children}
      </MantineButton>
    );
  }
);

Button.displayName = "StudioButton";
