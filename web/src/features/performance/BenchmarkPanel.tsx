import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { api } from '@/lib/api'
import { useUIStore } from '@/stores/uiStore'
import type {
  BenchmarkJob,
  BenchmarkModelSummary,
  BenchmarkWorkload,
  Model,
} from '@/types/api'
import {
  estimateBenchmarkMinutes,
  formatMs,
  formatRate,
  formatWhen,
  metricLabels,
  modelDisplayName,
} from './performanceFormat'

function SelectCard({
  selected,
  title,
  subtitle,
  meta,
  onClick,
}: {
  selected: boolean
  title: string
  subtitle?: string
  meta?: string
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-pressed={selected}
      className={[
        'min-w-0 rounded-xl p-4 text-left transition duration-150',
        selected
          ? 'bg-primary-soft shadow-[inset_0_0_0_1.5px_rgb(var(--color-primary)/0.55)]'
          : 'bg-surface shadow-[inset_0_0_0_1px_rgb(var(--color-line)/0.55)] hover:bg-raised/70',
      ].join(' ')}
    >
      <p
        className={[
          'font-semibold',
          selected ? 'text-primary-active' : 'text-ink',
        ].join(' ')}
      >
        {title}
      </p>
      {subtitle ? <p className="mt-1 text-sm text-ink-muted">{subtitle}</p> : null}
      {meta ? <p className="mt-2 text-xs text-ink-faint">{meta}</p> : null}
    </button>
  )
}

function VisualBars({
  rows,
  unit,
  higherIsBetter,
}: {
  rows: { id: string; label: string; value: number; winner?: boolean }[]
  unit: string
  higherIsBetter: boolean
}) {
  const max = Math.max(...rows.map((r) => r.value), 0.001)
  const sorted = [...rows].sort((a, b) =>
    higherIsBetter ? b.value - a.value : a.value - b.value,
  )
  return (
    <ul className="space-y-2">
      {sorted.map((r) => {
        const pct = Math.max(4, Math.round((r.value / max) * 100))
        return (
          <li key={r.id} className="space-y-1">
            <div className="flex items-baseline justify-between gap-2 text-sm">
              <span className={r.winner ? 'font-semibold text-ink' : 'text-ink'}>
                {r.label}
                {r.winner ? (
                  <span className="badge-preferred ml-2">Fastest</span>
                ) : null}
              </span>
              <span className="shrink-0 tabular-nums text-ink-muted">
                {higherIsBetter ? formatRate(r.value) : formatMs(r.value)}
                {higherIsBetter ? ` ${unit}` : ''}
              </span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-raised">
              <div
                className={[
                  'h-full rounded-full transition-all',
                  r.winner ? 'bg-accent' : 'bg-primary/70',
                ].join(' ')}
                style={{ width: `${pct}%` }}
              />
            </div>
          </li>
        )
      })}
    </ul>
  )
}

function BenchmarkResults({
  job,
  workloads,
  models,
}: {
  job: BenchmarkJob
  workloads: BenchmarkWorkload[]
  models: Model[]
}) {
  const advanced = useUIStore((s) => s.advancedMode)
  const labels = metricLabels(advanced)
  const [showDetail, setShowDetail] = useState(false)

  const byWorkload = useMemo(() => {
    const map = new Map<string, BenchmarkModelSummary[]>()
    for (const row of job.summaries ?? []) {
      const list = map.get(row.workload_id) ?? []
      list.push(row)
      map.set(row.workload_id, list)
    }
    return [...map.entries()]
  }, [job.summaries])

  const workloadName = (id: string) => workloads.find((w) => w.id === id)?.name ?? id

  return (
    <section className="card space-y-5">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h2 className="section-title">Results</h2>
          <p className="mt-1 text-sm capitalize text-ink-muted">Status: {job.status}</p>
        </div>
        {(job.status === 'running' || job.status === 'pending') && (
          <div className="min-w-[220px] flex-1">
            <div className="mb-1 flex justify-between text-xs text-ink-muted">
              <span>{job.progress.message || job.progress.phase}</span>
              <span className="tabular-nums">{job.progress.percent}%</span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-raised">
              <div
                className="h-full bg-primary transition-all"
                style={{ width: `${Math.min(job.progress.percent, 100)}%` }}
              />
            </div>
          </div>
        )}
      </div>

      {job.error && (
        <div className="rounded-lg border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
          {job.error}
        </div>
      )}

      {Object.keys(job.winners ?? {}).length > 0 && (
        <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
          {Object.entries(job.winners).map(([workloadId, modelId]) => (
            <div
              key={workloadId}
              className="rounded-lg border border-accent/30 bg-accent-soft px-4 py-3"
            >
              <p className="label-caps text-accent">
                {job.status === 'completed' ? 'Winner' : 'Leading'} ·{' '}
                {workloadName(workloadId)}
              </p>
              <p className="mt-1 font-medium text-ink">
                {modelDisplayName(models, modelId)}
              </p>
            </div>
          ))}
        </div>
      )}

      {byWorkload.length > 0 && (
        <div className="space-y-6">
          {byWorkload.map(([workloadId, rows]) => (
            <div key={workloadId} className="space-y-4">
              <h3 className="font-display text-base font-semibold text-ink">
                {workloadName(workloadId)}
              </h3>
              <div>
                <p className="mb-2 text-xs font-medium text-ink-muted">{labels.speed}</p>
                <VisualBars
                  higherIsBetter
                  unit="tok/s"
                  rows={rows.map((r) => ({
                    id: r.model_id,
                    label: modelDisplayName(models, r.model_id),
                    value: r.avg_eval_tok_per_sec,
                    winner: job.winners?.[workloadId] === r.model_id,
                  }))}
                />
              </div>
              <div>
                <p className="mb-2 text-xs font-medium text-ink-muted">
                  {labels.firstResponse}
                </p>
                <VisualBars
                  higherIsBetter={false}
                  unit=""
                  rows={rows.map((r) => ({
                    id: r.model_id,
                    label: modelDisplayName(models, r.model_id),
                    value: r.avg_ttft_ms,
                    winner:
                      Math.min(...rows.map((x) => x.avg_ttft_ms || Infinity)) ===
                      r.avg_ttft_ms,
                  }))}
                />
              </div>
            </div>
          ))}
        </div>
      )}

      {byWorkload.length > 0 && (
        <div>
          <button
            type="button"
            className="text-sm text-primary hover:underline"
            onClick={() => setShowDetail((o) => !o)}
          >
            {showDetail ? 'Hide detailed results' : 'Detailed results'}
          </button>
          {showDetail && (
            <div className="mt-3 space-y-4">
              {byWorkload.map(([workloadId, rows]) => (
                <div
                  key={workloadId}
                  className="min-w-0 overflow-hidden rounded-xl border border-line"
                >
                  <div className="border-b border-line bg-raised/50 px-4 py-2">
                    <p className="text-sm font-medium text-ink">
                      {workloadName(workloadId)}
                    </p>
                  </div>
                  <div className="table-scroll">
                    <table className="min-w-full text-left text-sm">
                      <thead>
                        <tr className="border-b border-line text-ink-muted">
                          <th className="px-3 py-2 font-medium">Model</th>
                          <th className="px-3 py-2 text-right font-medium">
                            {labels.speed}
                          </th>
                          <th className="px-3 py-2 text-right font-medium">
                            {labels.prompt}
                          </th>
                          <th className="px-3 py-2 text-right font-medium">
                            {labels.firstResponse}
                          </th>
                          <th className="px-3 py-2 text-right font-medium">
                            {labels.total}
                          </th>
                          {advanced && (
                            <th className="px-3 py-2 text-right font-medium">Load</th>
                          )}
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-line">
                        {rows.map((row) => (
                          <tr key={`${row.model_id}-${row.workload_id}`}>
                            <td className="px-3 py-2.5 text-ink">
                              {modelDisplayName(models, row.model_id)}
                            </td>
                            <td className="px-3 py-2.5 text-right tabular-nums">
                              {formatRate(row.avg_eval_tok_per_sec)}
                            </td>
                            <td className="px-3 py-2.5 text-right tabular-nums">
                              {formatRate(row.avg_prompt_tok_per_sec)}
                            </td>
                            <td className="px-3 py-2.5 text-right tabular-nums">
                              {formatMs(row.avg_ttft_ms)}
                            </td>
                            <td className="px-3 py-2.5 text-right tabular-nums">
                              {formatMs(row.avg_total_ms)}
                            </td>
                            {advanced && (
                              <td className="px-3 py-2.5 text-right tabular-nums">
                                {formatMs(row.load_ms ?? 0)}
                              </td>
                            )}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {(job.status === 'running' || job.status === 'pending') && byWorkload.length === 0 && (
        <LoadingSpinner label={job.progress.message || 'Running benchmark…'} />
      )}
    </section>
  )
}

export function BenchmarkPanel() {
  const queryClient = useQueryClient()
  const advanced = useUIStore((s) => s.advancedMode)
  const [selectedModels, setSelectedModels] = useState<string[]>([])
  const [selectedWorkloads, setSelectedWorkloads] = useState<string[]>([])
  const [runsPerPrompt, setRunsPerPrompt] = useState(2)
  const [activeJobId, setActiveJobId] = useState<string | null>(null)
  const [benchError, setBenchError] = useState<string | null>(null)

  const modelsQuery = useQuery({
    queryKey: ['models'],
    queryFn: () => api.getModels(),
    retry: false,
  })

  const workloadsQuery = useQuery({
    queryKey: ['benchmark-workloads'],
    queryFn: () => api.listBenchmarkWorkloads(),
    retry: false,
    staleTime: 60_000,
  })

  const jobQuery = useQuery({
    queryKey: ['benchmark', activeJobId],
    queryFn: () => (activeJobId ? api.getBenchmark(activeJobId) : null),
    enabled: Boolean(activeJobId),
    refetchInterval: (query) => {
      const status = query.state.data?.status
      if (status === 'running' || status === 'pending') return 1200
      return false
    },
    retry: false,
  })

  const jobsQuery = useQuery({
    queryKey: ['benchmarks'],
    queryFn: () => api.listBenchmarks(),
    retry: false,
    refetchInterval: (query) => {
      const items = query.state.data ?? []
      return items.some((j) => j.status === 'running' || j.status === 'pending')
        ? 2000
        : false
    },
  })

  const models = modelsQuery.data ?? []
  const installedModels = useMemo(
    () =>
      models.filter(
        (m) => m.installed || (m.installed_on?.length ?? 0) > 0,
      ),
    [models],
  )
  const workloads = workloadsQuery.data ?? []
  const activeJob = jobQuery.data ?? null
  const recentJobs = jobsQuery.data ?? []

  const promptCount = useMemo(() => {
    return workloads
      .filter((w) => selectedWorkloads.includes(w.id))
      .reduce((n, w) => n + w.prompts.length, 0)
  }, [workloads, selectedWorkloads])

  const estMinutes = estimateBenchmarkMinutes(
    selectedModels.length,
    promptCount,
    runsPerPrompt,
  )

  const startMutation = useMutation({
    mutationFn: () =>
      api.startBenchmark({
        model_ids: selectedModels,
        workload_ids: selectedWorkloads,
        runs: runsPerPrompt,
      }),
    onSuccess: (job) => {
      if (job?.id) {
        setActiveJobId(job.id)
        setBenchError(null)
        queryClient.invalidateQueries({ queryKey: ['benchmark', job.id] })
        queryClient.invalidateQueries({ queryKey: ['benchmarks'] })
      }
    },
    onError: (error) => {
      setBenchError(error instanceof Error ? error.message : 'Could not start benchmark.')
    },
  })

  const cancelMutation = useMutation({
    mutationFn: () =>
      activeJobId ? api.cancelBenchmark(activeJobId) : Promise.resolve(null),
    onSuccess: () => {
      if (activeJobId) {
        queryClient.invalidateQueries({ queryKey: ['benchmark', activeJobId] })
        queryClient.invalidateQueries({ queryKey: ['benchmarks'] })
      }
    },
  })

  const canStart =
    selectedModels.length >= 1 &&
    selectedWorkloads.length >= 1 &&
    !startMutation.isPending &&
    activeJob?.status !== 'running' &&
    activeJob?.status !== 'pending'

  const toggleModel = (id: string) => {
    setSelectedModels((current) =>
      current.includes(id) ? current.filter((x) => x !== id) : [...current, id].slice(0, 6),
    )
  }

  const toggleWorkload = (id: string) => {
    setSelectedWorkloads((current) =>
      current.includes(id) ? current.filter((x) => x !== id) : [...current, id],
    )
  }

  return (
    <div className="space-y-6 animate-fade">
      <section className="card space-y-6">
        <div>
          <h2 className="section-title">Compare models</h2>
          <p className="mt-1 text-sm text-ink-muted">
            Run the same prompts across models and see which is fastest for each kind of work.
          </p>
          <ol className="mt-4 flex flex-wrap gap-3 text-sm">
            {[
              { n: 1, label: 'Choose models', done: selectedModels.length > 0 },
              { n: 2, label: 'Choose workloads', done: selectedWorkloads.length > 0 },
              { n: 3, label: 'Run benchmark', done: false },
            ].map((step) => (
              <li
                key={step.n}
                className={[
                  'inline-flex items-center gap-2 rounded-full px-3 py-1',
                  step.done ? 'bg-success/15 text-success' : 'bg-raised text-ink-muted',
                ].join(' ')}
              >
                <span className="font-semibold tabular-nums">{step.n}</span>
                {step.label}
              </li>
            ))}
          </ol>
        </div>

        <div>
          <p className="label-caps mb-2">1. Choose models</p>
          {modelsQuery.isLoading && <LoadingSpinner label="Loading models…" />}
          {!modelsQuery.isLoading && installedModels.length === 0 && (
            <p className="text-sm text-ink-muted">
              Install at least one model on the Models page before benchmarking.
            </p>
          )}
          <div className="grid gap-2 sm:grid-cols-2">
            {installedModels.map((model) => (
              <SelectCard
                key={model.id}
                selected={selectedModels.includes(model.id)}
                title={model.display_name}
                subtitle={advanced ? model.id : undefined}
                onClick={() => toggleModel(model.id)}
              />
            ))}
          </div>
          <p className="mt-2 text-xs text-ink-muted">{selectedModels.length}/6 selected</p>
        </div>

        <div>
          <p className="label-caps mb-2">2. Choose workloads</p>
          <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
            {workloads.map((workload) => (
              <SelectCard
                key={workload.id}
                selected={selectedWorkloads.includes(workload.id)}
                title={workload.name}
                subtitle={workload.description}
                meta={`${workload.prompts.length} prompts`}
                onClick={() => toggleWorkload(workload.id)}
              />
            ))}
          </div>
        </div>

        <div className="space-y-3 border-t border-line/60 pt-5">
          <p className="label-caps">3. Run benchmark</p>
          <div className="rounded-lg border border-warning/30 bg-warning/10 px-4 py-3 text-sm text-ink">
            Yggdrasil will temporarily load and unload models while benchmarking. Active chats
            may be slower during the test.
          </div>
          <div className="flex flex-wrap items-end gap-4">
            {advanced && (
              <label className="text-sm text-ink-muted">
                Measured runs per prompt
                <select
                  className="field mt-1 block"
                  value={runsPerPrompt}
                  onChange={(e) => setRunsPerPrompt(Number(e.target.value))}
                >
                  {[1, 2, 3].map((n) => (
                    <option key={n} value={n}>
                      {n}
                    </option>
                  ))}
                </select>
              </label>
            )}
            <button
              type="button"
              className="btn-primary"
              disabled={!canStart}
              onClick={() => startMutation.mutate()}
            >
              {startMutation.isPending ? 'Starting…' : 'Run benchmark'}
            </button>
            {activeJob &&
              (activeJob.status === 'running' || activeJob.status === 'pending') && (
                <button
                  type="button"
                  className="btn-secondary"
                  disabled={cancelMutation.isPending}
                  onClick={() => cancelMutation.mutate()}
                >
                  Cancel
                </button>
              )}
          </div>
          {selectedModels.length > 0 && selectedWorkloads.length > 0 && (
            <p className="text-xs text-ink-faint">
              {selectedModels.length}{' '}
              {selectedModels.length === 1 ? 'model' : 'models'} ·{' '}
              {selectedWorkloads.length}{' '}
              {selectedWorkloads.length === 1 ? 'workload' : 'workloads'} · ~{estMinutes}{' '}
              {estMinutes === 1 ? 'minute' : 'minutes'}
            </p>
          )}
        </div>

        {benchError && (
          <div className="rounded-lg border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
            {benchError}
          </div>
        )}
      </section>

      {activeJob && (
        <BenchmarkResults job={activeJob} workloads={workloads} models={models} />
      )}

      {recentJobs.length > 0 && (
        <section className="card space-y-3">
          <div>
            <h2 className="section-title">Recent benchmarks</h2>
            <p className="mt-1 text-sm text-ink-muted">
              Jobs stay available while Yggdrasil is running.
            </p>
          </div>
          <ul className="divide-y divide-line rounded-xl border border-line">
            {recentJobs.slice(0, 8).map((job) => {
              const winnerCount = Object.keys(job.winners ?? {}).length
              const active = job.id === activeJobId
              return (
                <li key={job.id}>
                  <button
                    type="button"
                    onClick={() => setActiveJobId(job.id)}
                    className={[
                      'flex w-full items-center justify-between gap-3 px-4 py-3 text-left transition hover:bg-raised/50',
                      active ? 'bg-primary-soft/40' : '',
                    ].join(' ')}
                  >
                    <div>
                      <p className="font-medium capitalize text-ink">{job.status}</p>
                      <p className="mt-0.5 text-xs text-ink-muted">
                        {job.request.model_ids.length} models ·{' '}
                        {job.request.workload_ids.length} workloads ·{' '}
                        {formatWhen(job.created_at)}
                      </p>
                    </div>
                    <div className="text-right text-xs text-ink-muted">
                      {winnerCount > 0 ? (
                        <span className="text-accent">{winnerCount} winners</span>
                      ) : (
                        <span className="tabular-nums">{job.progress.percent}%</span>
                      )}
                    </div>
                  </button>
                </li>
              )
            })}
          </ul>
        </section>
      )}
    </div>
  )
}
