import { useEffect, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { subscribeEvents } from '@/lib/events'

type RuntimeNotify = {
  InitializeNotifications?: () => Promise<void>
  RequestNotificationAuthorization?: () => Promise<boolean>
  CheckNotificationAuthorization?: () => Promise<boolean>
  IsNotificationAvailable?: () => Promise<boolean>
  SendNotification?: (opts: {
    title: string
    body?: string
    identifier?: string
  }) => Promise<void>
}

function runtimeNotify(): RuntimeNotify | undefined {
  return (window as unknown as { runtime?: RuntimeNotify }).runtime
}

async function ensurePermission(): Promise<boolean> {
  const rt = runtimeNotify()
  try {
    if (rt?.IsNotificationAvailable) {
      const available = await rt.IsNotificationAvailable()
      if (!available) return false
    }
    if (rt?.InitializeNotifications) {
      await rt.InitializeNotifications()
    }
    if (rt?.CheckNotificationAuthorization) {
      const ok = await rt.CheckNotificationAuthorization()
      if (ok) return true
    }
    if (rt?.RequestNotificationAuthorization) {
      return await rt.RequestNotificationAuthorization()
    }
  } catch {
    // fall through to browser API
  }
  if (typeof Notification === 'undefined') return false
  if (Notification.permission === 'granted') return true
  if (Notification.permission === 'denied') return false
  const result = await Notification.requestPermission()
  return result === 'granted'
}

async function sendNotice(title: string, body: string) {
  const rt = runtimeNotify()
  try {
    if (rt?.SendNotification) {
      await rt.SendNotification({
        title,
        body,
        identifier: `ygg-${Date.now()}`,
      })
      return
    }
  } catch {
    // fall through
  }
  if (typeof Notification !== 'undefined' && Notification.permission === 'granted') {
    new Notification(title, { body })
  }
}

/** Listens for task/node events and shows OS notifications when Settings allow. */
export function NotificationHost() {
  const settingsQuery = useQuery({
    queryKey: ['settings'],
    queryFn: () => api.getSettings(),
    retry: false,
    staleTime: 30_000,
    refetchInterval: 60_000,
  })
  const ready = useRef(false)

  useEffect(() => {
    const notifyTask = settingsQuery.data?.notify_task_finish ?? true
    const notifyPeer = settingsQuery.data?.notify_peer_offline ?? true
    if (!notifyTask && !notifyPeer) return

    let cancelled = false
    void (async () => {
      const ok = await ensurePermission()
      if (!cancelled) ready.current = ok
    })()

    const unsubscribe = subscribeEvents({
      onEvent: (event) => {
        if (!ready.current) return
        if (notifyTask && event.type === 'task.completed') {
          void sendNotice('Task finished', 'Yggdrasil finished a background task.')
        }
        if (notifyTask && event.type === 'task.failed') {
          const msg =
            (event.payload?.error as string | undefined) ||
            'A background task failed.'
          void sendNotice('Task failed', msg)
        }
        if (notifyPeer && event.type === 'node.offline') {
          const name =
            (event.payload?.name as string | undefined) ||
            (event.payload?.node_name as string | undefined) ||
            'A paired computer'
          void sendNotice('Computer offline', `${name} is no longer responding.`)
        }
      },
    })

    return () => {
      cancelled = true
      unsubscribe()
    }
  }, [
    settingsQuery.data?.notify_task_finish,
    settingsQuery.data?.notify_peer_offline,
  ])

  return null
}
