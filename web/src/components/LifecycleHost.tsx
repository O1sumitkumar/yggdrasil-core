import { useEffect, useState, type ReactNode } from 'react'
import { onDesktopLifecycle, type DesktopLifecycleEvent } from '@/lib/desktopBridge'

/**
 * Listens for Wails shell lifecycle events (quit / background) and shows
 * a blocking overlay while the local service is being stopped.
 */
export function LifecycleHost({ children }: { children: ReactNode }) {
  const [event, setEvent] = useState<DesktopLifecycleEvent | null>(null)

  useEffect(() => onDesktopLifecycle(setEvent), [])

  return (
    <div className="h-full min-h-0 min-w-0 overflow-hidden">
      {children}
      {event?.state === 'stopping' ? (
        <div
          className="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-canvas/92 px-6 backdrop-blur-sm"
          role="status"
          aria-live="assertive"
          aria-busy="true"
        >
          <img
            src="/yggdrasil-mark.png"
            alt=""
            width={56}
            height={56}
            className="h-14 w-14 object-contain animate-boot-pulse"
            decoding="async"
          />
          <p className="mt-5 font-display text-xl font-semibold tracking-tight text-ink">
            Closing Yggdrasil
          </p>
          <p className="mt-2 max-w-sm text-center text-sm text-ink-muted">
            {event.message || 'Stopping local AI service…'}
          </p>
          <span
            className="mt-6 inline-block h-5 w-5 animate-spin rounded-full border-2 border-line border-t-primary"
            aria-hidden
          />
        </div>
      ) : null}
    </div>
  )
}
