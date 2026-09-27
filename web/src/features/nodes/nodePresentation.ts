import type { Accelerator, HardwareInventory, Model, Node } from '@/types/api'
import { purposeChips } from '@/features/models/modelPresentation'

function isMeaningful(s?: string | null): s is string {
  if (!s) return false
  const t = s.trim()
  if (!t || t === '.' || t === '/' || t === '-' || t === '—') return false
  if (t === '/' || t.includes('/')) {
    const parts = t.split('/').map((p) => p.trim())
    if (parts.every((p) => !p || p === '.')) return false
  }
  return true
}

function gbLabel(bytes?: number): string | null {
  if (bytes == null || bytes <= 0) return null
  const gb = bytes / 1024 ** 3
  if (gb >= 10) return `${Math.round(gb)} GB`
  return `${gb.toFixed(1)} GB`
}

function accelLine(accel?: Accelerator): string | null {
  if (!accel) return null
  const name = [accel.model, accel.vendor].find(isMeaningful)
  if (!name) return null
  const vram =
    gbLabel(accel.dedicated_vram_bytes) ||
    (accel.unified_memory_bytes
      ? `${gbLabel(accel.unified_memory_bytes)} unified memory`
      : null)
  if (vram && accel.unified_memory_bytes) return `${name} · ${vram}`
  if (vram && accel.dedicated_vram_bytes) return `${name} · ${vram} VRAM`
  return name
}

export type NodeHardwareView = {
  primary: string | null
  secondary: string | null
  memoryBytes: number
  unavailable: boolean
}

export function describeNodeHardware(hw?: HardwareInventory | null): NodeHardwareView {
  if (!hw) {
    return { primary: null, secondary: null, memoryBytes: 0, unavailable: true }
  }

  const cpu = isMeaningful(hw.cpu?.model) ? hw.cpu.model.trim() : null
  const accel = hw.accelerators?.[0]
  const accelText = accelLine(accel)
  const ram = gbLabel(hw.memory?.total_bytes)
  const memBytes = hw.memory?.total_bytes ?? 0
  const unified = accel?.unified_memory_bytes
    ? gbLabel(accel.unified_memory_bytes)
    : null

  // Apple-style: chip as primary, unified memory as secondary
  if (cpu && /apple|m[0-9]/i.test(cpu) && (unified || ram)) {
    return {
      primary: cpu,
      secondary: unified ? `${unified} unified memory` : `${ram} memory`,
      memoryBytes: unified
        ? accel!.unified_memory_bytes!
        : memBytes,
      unavailable: false,
    }
  }

  if (cpu && accelText) {
    return {
      primary: cpu,
      secondary: accelText,
      memoryBytes: memBytes || accel?.dedicated_vram_bytes || 0,
      unavailable: false,
    }
  }

  if (accelText) {
    return {
      primary: accelText,
      secondary: ram ? `${ram} RAM` : null,
      memoryBytes: memBytes || accel?.dedicated_vram_bytes || 0,
      unavailable: false,
    }
  }

  if (cpu) {
    return {
      primary: cpu,
      secondary: ram ? `${ram} RAM` : null,
      memoryBytes: memBytes,
      unavailable: false,
    }
  }

  if (ram) {
    return {
      primary: `${ram} memory`,
      secondary: null,
      memoryBytes: memBytes,
      unavailable: false,
    }
  }

  return { primary: null, secondary: null, memoryBytes: 0, unavailable: true }
}

export type MembershipKind = 'local' | 'paired' | 'nearby' | 'offline'

export function membershipKind(node: Node): MembershipKind {
  if (node.is_local) return 'local'
  if (!node.paired) return 'nearby'
  if (node.status === 'offline') return 'offline'
  return 'paired'
}

export function membershipLabel(kind: MembershipKind): string {
  switch (kind) {
    case 'local':
      return 'This computer'
    case 'paired':
      return 'Paired'
    case 'nearby':
      return 'Available to pair'
    case 'offline':
      return 'Offline'
  }
}

export function onlineLabel(node: Node): string {
  if (node.is_local) return 'Online'
  if (node.status === 'online') return 'Online'
  if (node.status === 'offline') return 'Offline'
  return 'Checking…'
}

/** What this computer is good for, from installed models + hardware hints. */
export function availableForLabels(models: Model[], hw?: HardwareInventory | null): string[] {
  const set = new Set<string>()
  for (const m of models) {
    for (const chip of purposeChips(m)) set.add(chip)
    if (m.tags?.includes('large')) set.add('Large models')
  }

  const mem =
    hw?.memory?.total_bytes ??
    hw?.accelerators?.[0]?.unified_memory_bytes ??
    hw?.accelerators?.[0]?.dedicated_vram_bytes ??
    0
  const vram =
    hw?.accelerators?.[0]?.dedicated_vram_bytes ??
    hw?.accelerators?.[0]?.unified_memory_bytes ??
    0

  if (set.size === 0) {
    if (mem >= 24 * 1024 ** 3 || vram >= 16 * 1024 ** 3) set.add('Large models')
    if (vram >= 8 * 1024 ** 3) set.add('Vision')
    set.add('General')
    set.add('Coding')
  } else if (mem >= 48 * 1024 ** 3 || vram >= 20 * 1024 ** 3) {
    set.add('Large models')
  }

  const order = ['General', 'Coding', 'Reasoning', 'Vision', 'Tools', 'Fast', 'Large models']
  return order.filter((l) => set.has(l)).slice(0, 4)
}

export function teamBestAt(fleet: Node[], models: Model[]): string[] {
  const set = new Set<string>()
  for (const n of fleet) {
    const onNode = models.filter((m) =>
      (m.installed_on ?? []).some((p) => p.node_id === n.id),
    )
    for (const label of availableForLabels(onNode, n.hardware)) set.add(label)
  }
  const order = ['Coding', 'Vision', 'Reasoning', 'Large models', 'General', 'Tools', 'Fast']
  return order.filter((l) => set.has(l)).slice(0, 3)
}

export function combinedMemoryBytes(fleet: Node[]): number {
  let total = 0
  for (const n of fleet) {
    const view = describeNodeHardware(n.hardware)
    total += view.memoryBytes
  }
  return total
}
