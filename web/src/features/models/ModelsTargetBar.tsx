import { Link } from 'react-router-dom'
import { formatBytes } from '@/lib/format'
import type { HardwareInventory, Node, RunningModelView } from '@/types/api'

/** Page-level install / recommendation target. */
export type ModelsTarget = 'local' | 'all' | string

export function resolveInstallNodeId(target: ModelsTarget): string | undefined {
  if (target === 'local') return undefined
  return target
}

function hardwareLine(hw: HardwareInventory | null | undefined): string {
  if (!hw) return 'Hardware details unavailable'
  const parts: string[] = []
  const accel = hw.accelerators?.[0]
  if (accel?.model) parts.push(accel.model)
  else if (hw.cpu?.model) parts.push(hw.cpu.model)
  const mem =
    accel?.unified_memory_bytes ||
    accel?.dedicated_vram_bytes ||
    hw.memory?.total_bytes ||
    0
  if (mem > 0) {
    const unified = Boolean(accel?.unified_memory_bytes)
    parts.push(`${formatBytes(mem)}${unified ? ' unified memory' : ' memory'}`)
  }
  return parts.join(' · ') || 'Hardware details unavailable'
}

function targetHardware(
  target: ModelsTarget,
  localHw: HardwareInventory | null,
  nodes: Node[],
): HardwareInventory | null {
  if (target === 'local' || target === 'all') return localHw
  return nodes.find((n) => n.id === target)?.hardware ?? null
}

function targetName(
  target: ModelsTarget,
  localHw: HardwareInventory | null,
  nodes: Node[],
): string {
  if (target === 'all') return 'All computers'
  if (target === 'local') return localHw?.hostname || 'This computer'
  return nodes.find((n) => n.id === target)?.name || 'Computer'
}

export function ModelsTargetBar({
  target,
  onTargetChange,
  localHardware,
  nodes,
  running,
}: {
  target: ModelsTarget
  onTargetChange: (next: ModelsTarget) => void
  localHardware: HardwareInventory | null
  nodes: Node[]
  running: RunningModelView[]
}) {
  const localNode = nodes.find((n) => n.is_local)
  const onlinePeers = nodes.filter(
    (n) => !n.is_local && n.paired && n.status === 'online',
  )
  const hasCluster = onlinePeers.length > 0 || nodes.filter((n) => n.paired).length > 1

  const hw = targetHardware(target, localHardware, nodes)
  const name = targetName(target, localHardware, nodes)

  const runningForTarget =
    target === 'all'
      ? running
      : target === 'local'
        ? running.filter(
            (r) =>
              !localNode ||
              r.node_id === localNode.id ||
              r.node_id === '' ||
              r.node_name === localHardware?.hostname,
          )
        : running.filter((r) => r.node_id === target)

  const line = hardwareLine(hw)
  const runningLabel = `${runningForTarget.length} model${
    runningForTarget.length === 1 ? '' : 's'
  } running`

  return (
    <div className="min-w-0 rounded-xl border border-line/70 bg-surface px-4 py-3 shadow-panel">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div className="min-w-0">
          <label className="flex flex-wrap items-center gap-2">
            <span className="text-sm font-medium text-ink-muted">Models for:</span>
            <select
              className="field max-w-full py-1.5 text-sm font-semibold text-ink"
              value={target}
              onChange={(e) => onTargetChange(e.target.value as ModelsTarget)}
              aria-label="Computer for model recommendations"
            >
              <option value="local">
                {localHardware?.hostname || localNode?.name || 'This computer'}
              </option>
              {hasCluster && <option value="all">All computers</option>}
              {onlinePeers.map((n) => (
                <option key={n.id} value={n.id}>
                  {n.name}
                </option>
              ))}
            </select>
          </label>
          <p className="mt-1.5 text-sm text-ink-muted">
            {target === 'all' ? (
              <>
                Cluster view · {runningLabel}
              </>
            ) : (
              <>
                {line}
                <span className="text-ink-faint"> · </span>
                {runningLabel}
              </>
            )}
          </p>
          <p className="mt-1 text-xs text-ink-faint">
            {target === 'all'
              ? 'Installs can go to every paired computer. Recommendations still emphasize what fits well.'
              : `Recommendations are based on ${name}’s hardware.`}
          </p>
        </div>
        <Link
          to="/nodes"
          className="shrink-0 text-sm font-medium text-primary underline-offset-2 hover:underline"
        >
          {hasCluster ? 'Manage computers' : 'Use another computer'}
        </Link>
      </div>
    </div>
  )
}
