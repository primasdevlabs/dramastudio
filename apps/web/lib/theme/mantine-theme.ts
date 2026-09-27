import { createTheme, MantineColorsTuple } from '@mantine/core'

const studioDark: MantineColorsTuple = [
  '#f1f5f9',
  '#cbd5e1',
  '#94a3b8',
  '#64748b',
  '#475569',
  '#334155',
  '#26334D',
  '#1A2335',
  '#111726',
  '#090D16',
]

const cyanAccent: MantineColorsTuple = [
  '#e0f7fa',
  '#b2ebf2',
  '#80deea',
  '#4dd0e1',
  '#26c6da',
  '#00bcd4',
  '#06B6D4',
  '#0097a7',
  '#00838f',
  '#006064',
]

export const studioTheme = createTheme({
  primaryColor: 'cyan',
  primaryShade: 6,
  colors: {
    dark: studioDark,
    cyan: cyanAccent,
  },
  fontFamily: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  defaultRadius: 'md',
  components: {
    Button: {
      defaultProps: {
        radius: 'md',
      },
    },
    Paper: {
      defaultProps: {
        radius: 'lg',
        bg: 'dark.8',
      },
    },
    Card: {
      defaultProps: {
        radius: 'lg',
        bg: 'dark.8',
        withBorder: true,
      },
    },
  },
})
