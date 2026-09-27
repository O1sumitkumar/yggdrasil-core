import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { api } from '@/lib/api'
import { describeNodeHardware } from '@/features/nodes/nodePresentation'
import type { Node, RunningModelView, Task } from '@/types/api'
import { formatRate, memoryUsePercent } from './performanceFormat'

function MemoryBar({ percent }: { percent: number }) {
  return (
    <div className="flex items-center gap-3">
      <span className="w-16 shrink-0 text-xs text-ink-muted">Memory</span>
      <div className="h-2 flex-1 overflow-hidden rounded-full bg-raised">
        <div
          className="h-full bg-primary/80 transition-all duration-300"
          style={{ width: `${percent}%` }}
        />
      </div>
      <span className="w-10 shrink-0 text-right text-xs tabular-nums text-ink">{percent}%</span>
    </div>
  )
}

function NodeLiveCard({
  node,
  running,
}: {
  node: Node
  running: RunningModelView[]
}) {
  const hw = describeNodeHardware(node.hardware)
  const memPct = memoryUsePercent(
    node.hardware?.memory?.total_bytes,
    node.hardware?.memory?.available_bytes,
  )
  const online = node.is_local || node.status === 'online'
  const idle = running.length === 0

  return (
    <article
      className={[
        'card space-y-3',
        !online ? 'opacity-70' : '',
      ].join(' ')}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0">
          <h3 className="font-display text-lg font-semibold text-ink">{node.name}</h3>
          {!hw.unavailable && (
            <p className="mt-0.5 text-sm text-ink-muted">
              {[hw.primary, hw.secondary].filter(Boolean).join(' · ')}
            </p>
          )}
        </div>
        <span
          className={[
            'status-chip shrink-0',
            online ? 'bg-success/15 text-success' : 'bg-raised text-ink-muted',
          ].join(' ')}
        >
          {online ? (idle ? 'Idle' : 'Active') : 'Offline'}
        </span>
      </div>

      {!online ? (
        <p className="text-sm text-ink-faint">Not reachable right now.</p>
      ) : idle ? (
        <p className="text-sm text-ink-muted">Idle — no models loaded</p>
      ) : (
        <div className="space-y-3">
          {memPct != null && <MemoryBar percent={memPct} />}
          <ul className="space-y-3">
            {running.map((r) => (
              <li key={r.instance_id} className="space-y-1 text-sm">
                <div className="flex justify-between gap-2">
                  <span className="text-ink-muted">Model</span>
                  <span className="truncate font-medium text-ink">{r.display_name}</span>
                </div>
                <div className="flex justify-between gap-2">
                  <span className="text-ink-muted">Speed</span>
                  <span className="tabular-nums text-ink">
                    {r.speed_tok_per_sec && r.speed_tok_per_sec > 0
                      ? `${formatRate(r.speed_tok_per_sec)} tok/s`
                      : '—'}
                  </span>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </article>
  )
}

export function OverviewPanel() {
  const nodesQuery = useQuery({
    queryKey: ['nodes'],
    queryFn: () => api.getNodes(),
    retry: false,
    refetchInterval: 5_000,
  })

  const runningQuery = useQuery({
    queryKey: ['models-running'],
    queryFn: () => api.listRunningModels(),
    retry: false,
    refetchInterval: 3_000,
  })

  const tasksQuery = useQuery({
    queryKey: ['tasks'],
    queryFn: () => api.listTasks(),
    retry: false,
    refetchInterval: 5_000,
  })

  const fleet = (nodesQuery.data ?? []).filter((n) => n.is_local || n.paired)
  const running = runningQuery.data ?? []
  const tasks = (tasksQuery.data ?? []).filter(
    (t: Task) => t.status === 'running' || t.status === 'pending',
  )

  const connected = fleet.filter((n) => n.is_local || n.status === 'online').length
  const crossing =
    running.length > 0 &&
    new Set(running.map((r) => r.node_id).filter(Boolean)).size > 1

  const byNode = new Map<string, RunningModelView[]>()
  for (const r of running) {
    const id = r.node_id || ''
    const list = byNode.get(id) ?? []
    list.push(r)
    byNode.set(id, list)
  }

  if (nodesQuery.isLoading) {
    return <LoadingSpinner label="Checking what your AI is doing…" />
  }

  return (
    <div className="space-y-6 animate-fade">
      <section className="rounded-2xl bg-raised/40 px-5 py-4">
        <p className="text-sm text-ink">
          {connected} {connected === 1 ? 'computer' : 'computers'} connected
          {' · '}
          {running.length} {running.length === 1 ? 'model' : 'models'} running
          {' · '}
          {tasks.length} active {tasks.length === 1 ? 'task' : 'tasks'}
          {crossing ? ' · work across machines' : ''}
        </p>
        {tasks.length > 0 && (
          <ul className="mt-3 space-y-1.5">
            {tasks.slice(0, 3).map((t) => (
              <li key={t.id} className="truncate text-xs text-ink-muted">
                <span className="font-medium text-accent capitalize">{t.status}</span>
                {' — '}
                {t.prompt?.slice(0, 80) || 'Task in progress'}
              </li>
            ))}
          </ul>
        )}
      </section>

      {fleet.length === 0 ? (
        <p className="text-sm text-ink-muted">
          No computers in your team yet.{' '}
          <Link to="/nodes" className="text-primary hover:underline">
            Add a computer
          </Link>
        </p>
      ) : (
        <ul className="grid gap-4 md:grid-cols-2">
          {fleet.map((node) => (
            <li key={node.id}>
              <NodeLiveCard
                node={node}
                running={byNode.get(node.id) ?? []}
              />
            </li>
          ))}
        </ul>
      )}

      {running.length === 0 && fleet.length > 0 && (
        <p className="text-sm text-ink-faint">
          Nothing loaded right now. Open{' '}
          <Link to="/chat" className="text-primary hover:underline">
            Chat
          </Link>{' '}
          or run a benchmark when you want to compare models.
        </p>
      )}
    </div>
  )
}
