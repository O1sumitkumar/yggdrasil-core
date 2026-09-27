import { useQuery } from '@tanstack/react-query'
import { useEffect, useState, type ReactNode } from 'react'
import { api } from '@/lib/api'
import { readScreenshotLaunch } from '@/lib/screenshotMode'

const BOOT_GIVE_UP_MS = 25_000

function BootSplash({ message }: { message: string }) {
  return (
    <div
      className="flex h-full min-h-0 flex-col items-center justify-center bg-canvas px-6"
      role="status"
      aria-live="polite"
      aria-busy="true"
    >
      <img
        src="/yggdrasil-mark.png"
        alt=""
        width={72}
        height={72}
        className="h-[72px] w-[72px] object-contain animate-boot-pulse"
        decoding="async"
      />
      <p className="mt-6 font-display text-xl font-semibold tracking-tight text-ink">
        Starting Yggdrasil
      </p>
      <p className="mt-2 text-sm text-ink-muted">{message}</p>
      <span
        className="mt-6 inline-block h-5 w-5 animate-spin rounded-full border-2 border-line border-t-primary"
        aria-hidden
      />
    </div>
  )
}

function BootFailed({ onRetry, busy }: { onRetry: () => void; busy: boolean }) {
  return (
    <div className="flex h-full min-h-0 flex-col items-center justify-center bg-canvas px-6 text-center">
      <img
        src="/yggdrasil-mark.png"
        alt=""
        width={64}
        height={64}
        className="h-16 w-16 object-contain opacity-80"
        decoding="async"
      />
      <h1 className="mt-6 font-display text-2xl font-semibold tracking-tight text-ink">
        Service unavailable
      </h1>
      <p className="mt-2 max-w-md text-sm leading-relaxed text-ink-muted">
        The local Yggdrasil service has not responded yet. Wait a moment and try again — if this
        keeps happening, check Diagnostics after the app loads.
      </p>
      <button
        type="button"
        className="btn-primary mt-6"
        disabled={busy}
        onClick={onRetry}
      >
        {busy ? 'Retrying…' : 'Try again'}
      </button>
    </div>
  )
}

/**
 * Blocks the main shell until the local daemon answers /api/v1/health.
 * Screenshot mode paints the shell immediately with the demo API.
 */
export function ServiceBootGate({ children }: { children: ReactNode }) {
  if (readScreenshotLaunch()?.enabled) {
    return <>{children}</>
  }
  return <DaemonBootGate>{children}</DaemonBootGate>
}

function DaemonBootGate({ children }: { children: ReactNode }) {
  const [deadline, setDeadline] = useState(() => Date.now() + BOOT_GIVE_UP_MS)
  const [timedOut, setTimedOut] = useState(false)

  const healthQuery = useQuery({
    queryKey: ['health'],
    queryFn: async () => {
      const health = await api.getHealth()
      if (!health || health.status !== 'ok') {
        throw new Error('Service not ready')
      }
      return health
    },
    retry: true,
    retryDelay: (attempt) => Math.min(400 + attempt * 250, 2000),
    refetchInterval: (query) =>
      query.state.data?.status === 'ok' ? 15_000 : 1_000,
    refetchOnWindowFocus: true,
  })

  const ready = healthQuery.data?.status === 'ok'

  useEffect(() => {
    if (ready) {
      setTimedOut(false)
      return
    }
    const id = window.setInterval(() => {
      if (Date.now() >= deadline) {
        setTimedOut(true)
      }
    }, 400)
    return () => window.clearInterval(id)
  }, [ready, deadline])

  if (ready) {
    return <div className="h-full min-h-0 min-w-0 overflow-hidden">{children}</div>
  }

  if (timedOut) {
    return (
      <BootFailed
        busy={healthQuery.isFetching}
        onRetry={() => {
          setTimedOut(false)
          setDeadline(Date.now() + BOOT_GIVE_UP_MS)
          void healthQuery.refetch()
        }}
      />
    )
  }

  return <BootSplash message="Preparing your local AI service" />
}
