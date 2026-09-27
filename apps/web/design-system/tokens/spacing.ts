/**
 * DramaStudio Design System Tokens - Spacing, Radius, Shadows & Motion
 */

export const spacing = {
  xs: "0.25rem",
  sm: "0.5rem",
  md: "1rem",
  lg: "1.5rem",
  xl: "2rem",
  "2xl": "3rem",
} as const;

export const radius = {
  sm: "0.375rem",
  md: "0.5rem",
  lg: "0.75rem",
  xl: "1rem",
  full: "9999px",
} as const;

export const shadows = {
  card: "0 10px 30px -10px rgba(0, 0, 0, 0.5)",
  panel: "0 2px 8px rgba(0, 0, 0, 0.4)",
  dropdown: "0 4px 16px rgba(0, 0, 0, 0.5)",
} as const;

export const motion = {
  transitionFast: "all 150ms cubic-bezier(0.4, 0, 0.2, 1)",
  transitionNormal: "all 250ms cubic-bezier(0.4, 0, 0.2, 1)",
} as const;

