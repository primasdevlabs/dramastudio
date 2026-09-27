"use client";

import { Button as MantineButton, ButtonProps as MantineButtonProps } from "@mantine/core";
import { forwardRef } from "react";

export interface StudioButtonProps extends MantineButtonProps {
  variant?: "primary" | "secondary" | "subtle" | "destructive";
}

export const Button = forwardRef<HTMLButtonElement, StudioButtonProps>(
  ({ variant = "primary", className = "", children, ...props }, ref) => {
    let color = "terracotta";
    let mantineVariant: MantineButtonProps["variant"] = "filled";

    if (variant === "secondary") {
      mantineVariant = "outline";
      color = "gray";
    } else if (variant === "subtle") {
      mantineVariant = "subtle";
      color = "gray";
    } else if (variant === "destructive") {
      mantineVariant = "filled";
      color = "red";
    }

    return (
      <MantineButton
        ref={ref}
        radius="sm"
        color={color}
        variant={mantineVariant}
        className={`font-semibold text-xs tracking-wide ${className}`}
        {...props}
      >
        {children}
      </MantineButton>
    );
  }
);

Button.displayName = "StudioButton";
