import './globals.css'
import { ColorSchemeScript, MantineProvider } from '@mantine/core'
import { Notifications } from '@mantine/notifications'
import { Inter } from 'next/font/google'
import { QueryProvider } from '@/providers/query-provider'
import { RealtimeProvider } from '@/providers/realtime-provider'
import { studioTheme } from '@/lib/theme/mantine-theme'

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
})

export const metadata = {
  title: 'DramaStudio - Digital Drama Production Studio',
  description: 'Digital drama production studio and Lead Director system',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" data-mantine-color-scheme="dark" className={inter.variable}>
      <head>
        <ColorSchemeScript defaultColorScheme="dark" />
      </head>
      <body className={`${inter.className} bg-studio-bg text-studio-text antialiased min-h-screen`}>
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
