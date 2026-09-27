import './globals.css'
import { QueryProvider } from '@/providers/query-provider'
import { RealtimeProvider } from '@/providers/realtime-provider'

export const metadata = {
  title: 'DramaStudio - AI Autonomous Production Studio',
  description: 'AI-native autonomous drama production studio',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-studio-bg text-studio-text antialiased min-h-screen">
        <QueryProvider>
          <RealtimeProvider>
            {children}
          </RealtimeProvider>
        </QueryProvider>
      </body>
    </html>
  )
}
