'use client'

import { createContext, useContext, useEffect, useState } from 'react'

interface RealtimeEvent {
  id: string
  type: string
  message: string
  timestamp: string
}

const RealtimeContext = createContext<{ events: RealtimeEvent[] }>({ events: [] })

export function RealtimeProvider({ children }: { children: React.ReactNode }) {
  const [events, setEvents] = useState<RealtimeEvent[]>([])

  useEffect(() => {
    const initialEvents: RealtimeEvent[] = [
      { id: '1', type: 'AgentTaskCreated', message: 'Lead Director initiated Season 1 Episode 12 storyboard workflow', timestamp: 'Just now' },
      { id: '2', type: 'GenerationStarted', message: 'Visual Agent started Wan video generation for Shot 07', timestamp: '1m ago' },
      { id: '3', type: 'ContinuityIssueDetected', message: 'Continuity Agent flagged Wardrobe mismatch on Sarah (Scene 4)', timestamp: '3m ago' },
    ]
    setEvents(initialEvents)
  }, [])

  return (
    <RealtimeContext.Provider value={{ events }}>
      {children}
    </RealtimeContext.Provider>
  )
}

export const useRealtimeEvents = () => useContext(RealtimeContext)
