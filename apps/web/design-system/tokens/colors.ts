/**
 * DramaStudio Design System Tokens - Colors
 * Professional creative software palette (DaVinci Resolve / Linear / Figma inspired).
 * No gradients, no glowing effects, no purple/indigo AI aesthetic, no generic SaaS blue.
 */

export const colors = {
  // Dark Production Theme (Default for Workstation)
  dark: {
    canvas: "#151514",
    surface: "#1C1B19",
    surfaceRaised: "#24221F",
    ink: "#F2EEE7",
    inkMuted: "#A9A39A",
    border: "#393631",
    accent: "#C56A45", // Burnt Terracotta
    accentDark: "#8E3F25",
    success: "#6C9A76", // Muted Green
    warning: "#C18A45", // Warm Amber
    danger: "#C65A50",  // Muted Red
  },
  // Light Warm Neutral Theme
  light: {
    canvas: "#F4F1EB",
    surface: "#FFFCF7",
    surfaceRaised: "#FFFFFF",
    ink: "#191816",
    inkMuted: "#68645D",
    border: "#D9D4CB",
    borderStrong: "#BEB8AE",
    accent: "#B45732",
    accentDark: "#8E3F25",
    success: "#3E6B4A",
    warning: "#A66A24",
    danger: "#A83E36",
    info: "#4D6672",
  },
  // Production Semantics
  semantics: {
    completed: "#6C9A76",
    running: "#C56A45",
    waiting: "#A9A39A",
    warning: "#C18A45",
    blocked: "#C65A50",
    draft: "#A9A39A",
    approved: "#6C9A76",
    rejected: "#C65A50",
  },
} as const;
