import { getApiBase } from '@/lib/api'
import type { YggdrasilEvent } from '@/types/api'

export interface EventSubscriptionOptions {
  onEvent: (event: YggdrasilEvent) => void
  onError?: (error: Event) => void
  onOpen?: () => void
}

function parseEventData(raw: string): YggdrasilEvent | null {
  try {
    return JSON.parse(raw) as YggdrasilEvent
  } catch {
    return null
  }
}

const KNOWN_EVENT_TYPES = [
  'chat.token',
  'chat.complete',
  'chat.error',
  'model.download.started',
  'model.download.progress',
  'model.download.completed',
  'model.download.failed',
  'model.load.started',
  'model.load.completed',
  'orchestration.role',
  'orchestration.final',
  'scheduler.placement',
  'agent.started',
  'agent.completed',
  'tool.requested',
  'tool.started',
  'tool.completed',
  'tool.failed',
  'tool.parsed',
  'chat.model_routed',
  'task.created',
  'task.started',
  'task.completed',
  'task.failed',
  'node.discovered',
  'node.paired',
  'node.online',
  'node.offline',
] as const

export function subscribeEvents({
  onEvent,
  onError,
  onOpen,
}: EventSubscriptionOptions): () => void {
  // jsdom / older runtimes may not provide EventSource.
  if (typeof EventSource === 'undefined') {
    return () => {}
  }

  const url = `${getApiBase()}/api/v1/events`
  const source = new EventSource(url)

  source.onopen = () => {
    onOpen?.()
  }

  source.onmessage = (message) => {
    const event = parseEventData(message.data)
    if (event) {
      onEvent(event)
    }
  }

  for (const type of KNOWN_EVENT_TYPES) {
    source.addEventListener(type, (message) => {
      const eventMessage = message as MessageEvent<string>
      const event = parseEventData(eventMessage.data)
      if (event) {
        onEvent(event)
      }
    })
  }

  source.onerror = (error) => {
    onError?.(error)
  }

  return () => {
    source.close()
  }
}
