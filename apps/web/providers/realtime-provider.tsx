'use client'

import { createContext, useContext, useEffect, useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9471'
const MAX_EVENTS = 50

export interface RealtimeEvent {
  id: string
  type: string
  message: string
  timestamp: string
  project_id?: string
  aggregate_id?: string
}

interface DomainEventPayload {
  type: string
  version?: number
  aggregate?: string
  aggregate_id?: string
  project_id?: string
  message?: string
  occurred_at?: string
}

const RealtimeContext = createContext<{ events: RealtimeEvent[]; connected: boolean }>({
  events: [],
  connected: false,
})

/** Invalidates React Query caches that a domain event may have changed. */
function queriesForEvent(type: string, projectId?: string): string[][] {
  const keys: string[][] = []
  if (projectId) keys.push(['projects', projectId])
  switch (type) {
    case 'ProjectCreated':
    case 'ProjectUpdated':
      keys.push(['projects'])
      break
    case 'EpisodeCreated':
    case 'EpisodeCompleted':
    case 'SceneCreated':
      keys.push(['episodes'], ['scenes'])
      break
    case 'StoryFactEstablished':
    case 'StoryFactChanged':
      keys.push(['canon'])
      break
    case 'CharacterCreated':
    case 'CharacterUpdated':
    case 'CharacterLocked':
      keys.push(['characters'])
      break
    case 'LocationAdded':
      keys.push(['world'], ['locations'])
      break
    case 'ContinuityIssueDetected':
    case 'ContinuityIssueResolved':
      keys.push(['continuity'])
      break
    case 'AssetGenerated':
    case 'AssetApproved':
    case 'AssetRejected':
    case 'GenerationStarted':
    case 'GenerationCompleted':
      keys.push(['assets'], ['media-jobs'])
      break
    case 'ProductionJobCreated':
    case 'ProductionJobCompleted':
    case 'ProductionJobFailed':
    case 'ProductionRunStarted':
    case 'ProductionRunPaused':
    case 'ProductionRunResumed':
    case 'ProductionRunStopped':
      keys.push(['production'])
      break
    case 'ApprovalRequested':
    case 'ApprovalGranted':
    case 'ApprovalRejected':
      keys.push(['production'], ['approvals'])
      break
    case 'AgentTaskCreated':
    case 'AgentTaskCompleted':
    case 'AgentTaskFailed':
    case 'AgentDecisionMade':
      keys.push(['agents'])
      break
    case 'RenderCompleted':
    case 'RenderFailed':
      keys.push(['renders'], ['postproduction'])
      break
    case 'PublicationPublished':
    case 'PublicationFailed':
      keys.push(['publishing'])
      break
  }
  return keys
}

export function RealtimeProvider({ children }: { children: React.ReactNode }) {
  const [events, setEvents] = useState<RealtimeEvent[]>([])
  const [connected, setConnected] = useState(false)
  const queryClient = useQueryClient()
  const counter = useRef(0)

  useEffect(() => {
    const source = new EventSource(`${API_BASE_URL}/v1/events`)

    source.onopen = () => setConnected(true)
    source.onerror = () => setConnected(false)

    // Server emits `event: <Type>` + `data: <DomainEvent JSON>`.
    const onEvent = (raw: MessageEvent) => {
      let parsed: DomainEventPayload
      try {
        parsed = JSON.parse(raw.data)
      } catch {
        return
      }
      counter.current += 1
      setEvents((prev) => [
        {
          id: `${parsed.occurred_at ?? Date.now()}-${counter.current}`,
          type: parsed.type ?? raw.type,
          message: parsed.message ?? '',
          timestamp: parsed.occurred_at ?? new Date().toISOString(),
          project_id: parsed.project_id,
          aggregate_id: parsed.aggregate_id,
        },
        ...prev.slice(0, MAX_EVENTS - 1),
      ])
      for (const key of queriesForEvent(parsed.type ?? '', parsed.project_id)) {
        queryClient.invalidateQueries({ queryKey: key })
      }
    }

    // Named SSE events need a listener per type; the bus vocabulary is
    // finite, so register each known type plus a fallback for unknown ones.
    const types = [
      'ProjectCreated', 'ProjectUpdated', 'SeriesCreated', 'SeasonCreated', 'ArcCreated',
      'EpisodeCreated', 'EpisodeApproved', 'EpisodeCompleted', 'SceneCreated', 'SceneApproved',
      'CharacterCreated', 'CharacterUpdated', 'CharacterLocked', 'LocationAdded',
      'StoryFactEstablished', 'StoryFactChanged', 'PlotThreadCreated',
      'GenerationStarted', 'GenerationCompleted', 'AssetGenerated', 'AssetApproved', 'AssetRejected',
      'ProductionJobCreated', 'ProductionJobCompleted', 'ProductionJobFailed',
      'ProductionRunStarted', 'ProductionRunPaused', 'ProductionRunResumed', 'ProductionRunStopped',
      'ContinuityIssueDetected', 'ContinuityIssueResolved',
      'AgentTaskCreated', 'AgentTaskCompleted', 'AgentTaskFailed', 'AgentDecisionMade',
      'ApprovalRequested', 'ApprovalGranted', 'ApprovalRejected',
      'RenderCompleted', 'RenderFailed', 'PublicationPublished', 'PublicationFailed',
      'SeriesBibleUpdated',
    ]
    for (const t of types) {
      source.addEventListener(t, onEvent)
    }

    return () => {
      setConnected(false)
      source.close()
    }
  }, [queryClient])

  return (
    <RealtimeContext.Provider value={{ events, connected }}>
      {children}
    </RealtimeContext.Provider>
  )
}

export const useRealtimeEvents = () => useContext(RealtimeContext)
