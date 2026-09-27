import './globals.css'
import { ColorSchemeScript, MantineProvider } from '@mantine/core'
import { Notifications } from '@mantine/notifications'
import { QueryProvider } from '@/providers/query-provider'
import { RealtimeProvider } from '@/providers/realtime-provider'
import { studioTheme } from '@/lib/theme/mantine-theme'

export const metadata = {
  title: 'DramaStudio - AI Autonomous Production Studio',
  description: 'AI-native autonomous drama production studio',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" data-mantine-color-scheme="dark">
      <head>
        <ColorSchemeScript defaultColorScheme="dark" />
      </head>
      <body className="bg-studio-bg text-studio-text antialiased min-h-screen">
        <MantineProvider theme={studioTheme} defaultColorScheme="dark" forceColorScheme="dark">
          <Notifications position="top-right" />
          <QueryProvider>
            <RealtimeProvider>
              {children}
            </RealtimeProvider>
          </QueryProvider>
        </MantineProvider>
      </body>
    </html>
  )
}
