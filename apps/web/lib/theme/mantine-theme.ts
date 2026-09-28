import { createTheme, MantineColorsTuple } from '@mantine/core'

const studioDark: MantineColorsTuple = [
  '#f2eee7',
  '#d9d3c8',
  '#a9a39a',
  '#78736a',
  '#524e47',
  '#393631',
  '#2d2b27',
  '#24221f',
  '#1c1b19',
  '#151514',
]

const terracottaAccent: MantineColorsTuple = [
  '#fbf2ee',
  '#f5e0d7',
  '#ecc0ae',
  '#e29c81',
  '#da7d59',
  '#d4683f',
  '#C56A45', // Primary Terracotta
  '#8E3F25', // Dark Terracotta
  '#74341e',
  '#612c1b',
]

export const studioTheme = createTheme({
  primaryColor: 'terracotta',
  primaryShade: 6,
  colors: {
    dark: studioDark,
    terracotta: terracottaAccent,
  },
  fontFamily: "var(--font-inter), system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  fontFamilyMonospace: "var(--font-inter), system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  defaultRadius: 'sm',
  components: {
    Button: {
      defaultProps: {
        radius: 'sm',
      },
    },
    Paper: {
      defaultProps: {
        radius: 'md',
        bg: 'dark.8',
      },
    },
    Card: {
      defaultProps: {
        radius: 'md',
        bg: 'dark.8',
        withBorder: true,
      },
    },
  },
})
