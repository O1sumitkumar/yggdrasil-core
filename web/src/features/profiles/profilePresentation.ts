import type { AIProfile, Purpose, ToolPolicy } from '@/types/api'

/** Built-in profile template IDs (must match server presets). */
export const BUILT_IN_PROFILE_IDS = new Set([
  'general-assistant',
  'programming',
  'research',
  'custom',
])

export function isBuiltInProfile(id: string): boolean {
  return BUILT_IN_PROFILE_IDS.has(id)
}

export const PURPOSE_LABELS: Record<string, string> = {
  general: 'Everyday conversation',
  coding: 'Programming & code help',
  research: 'Research & long-form reading',
  custom: 'Custom setup',
}

export function orchestratorLabel(
  orchestratorId: string,
  advanced: boolean,
): { title: string; detail: string } {
  if (orchestratorId === 'team') {
    return {
      title: advanced ? 'AI team · Orchestrator: Team' : 'AI team',
      detail:
        'Splits work into coordinator → worker → reviewer. Great with more than one computer or a plan-and-review loop.',
    }
  }
  return {
    title: advanced ? 'Single assistant · Orchestrator: Simple' : 'Single assistant',
    detail:
      'One assistant replies in a single pass. Best for everyday chat, and the mode that can use local tools.',
  }
}

export function computerSelectionLabel(mode: string): {
  title: string
  short: string
  detail: string
} {
  switch (mode) {
    case 'prefer_local':
      return {
        title: 'Prefer this computer',
        short: 'This computer first',
        detail:
          'Tries to run on this machine first. Falls back to a paired computer when a model is not installed here.',
      }
    case 'manual':
      return {
        title: 'Custom',
        short: 'Custom',
        detail:
          'Uses role pins for placement. Ideal when you want a specific step on a specific machine.',
      }
    default:
      return {
        title: 'Automatic',
        short: 'Automatic',
        detail:
          'Places each step on whichever paired computer already has the model. No IP or port setup.',
      }
  }
}

export function toolSummary(tools: ToolPolicy[] | undefined): {
  enabled: number
  ask: number
  allow: number
  deny: number
} {
  const list = tools ?? []
  let ask = 0
  let allow = 0
  let deny = 0
  for (const t of list) {
    if (t.policy === 'deny') deny += 1
    else if (t.policy === 'ask' || t.policy === 'allow-for-session') ask += 1
    else allow += 1
  }
  return { enabled: list.length - deny, ask, allow, deny }
}

/** Compact tool chips for the card face (first few interesting ones). */
export function toolChipPreview(tools: ToolPolicy[] | undefined): {
  label: string
  tone: 'ok' | 'ask' | 'off'
}[] {
  const byId = new Map((tools ?? []).map((t) => [t.tool_id, t.policy]))
  const rows: { id: string; label: string }[] = [
    { id: 'filesystem.read', label: 'Files' },
    { id: 'git.status', label: 'Git' },
    { id: 'terminal', label: 'Terminal' },
  ]
  const chips: { label: string; tone: 'ok' | 'ask' | 'off' }[] = []
  for (const row of rows) {
    const policy = byId.get(row.id)
    if (!policy) continue
    if (policy === 'deny') chips.push({ label: row.label, tone: 'off' })
    else if (policy === 'allow') chips.push({ label: `${row.label} ✓`, tone: 'ok' })
    else chips.push({ label: `${row.label} Ask`, tone: 'ask' })
  }
  return chips
}

export function roleDisplayName(role: string): string {
  if (!role) return 'Role'
  return role.charAt(0).toUpperCase() + role.slice(1)
}

export function roleHint(role: ModelRoleLike): string {
  const model = role.model_id ? shortModel(role.model_id) : 'Automatic model'
  const computer = role.node_id ? 'Pinned computer' : 'Automatic computer'
  return `${model} · ${computer}`
}

export function roleHelp(role: string): string {
  switch (role.toLowerCase()) {
    case 'coordinator':
      return 'Plans the approach before deeper work begins.'
    case 'worker':
      return 'Does the main drafting or coding step.'
    case 'reviewer':
      return 'Checks the draft and produces the final answer.'
    case 'researcher':
      return 'Focuses on careful reading and synthesis.'
    case 'assistant':
      return 'The single voice that talks with you.'
    default:
      return 'A named step in this profile’s workflow.'
  }
}

type ModelRoleLike = { role: string; model_id?: string; node_id?: string }

function shortModel(id: string): string {
  const base = id.split('/').pop() || id
  return base.length > 28 ? `${base.slice(0, 26)}…` : base
}

export function purposeIcon(purpose: string): 'general' | 'coding' | 'research' | 'custom' {
  if (purpose === 'coding' || purpose === 'research' || purpose === 'custom') {
    return purpose
  }
  return 'general'
}

export function sortProfilesForDisplay(profiles: AIProfile[]): AIProfile[] {
  const purposeOrder = ['general', 'coding', 'research', 'custom']
  return [...profiles].sort((a, b) => {
    const aBuilt = isBuiltInProfile(a.id) ? 0 : 1
    const bBuilt = isBuiltInProfile(b.id) ? 0 : 1
    if (aBuilt !== bBuilt) return aBuilt - bBuilt
    const ap = purposeOrder.indexOf(a.purpose)
    const bp = purposeOrder.indexOf(b.purpose)
    if (ap !== bp) return (ap === -1 ? 99 : ap) - (bp === -1 ? 99 : bp)
    return a.name.localeCompare(b.name)
  })
}

export type ProfileFilter = 'all' | 'builtin' | 'custom'

export function filterProfiles(
  profiles: AIProfile[],
  filter: ProfileFilter,
): AIProfile[] {
  if (filter === 'builtin') return profiles.filter((p) => isBuiltInProfile(p.id))
  if (filter === 'custom') return profiles.filter((p) => !isBuiltInProfile(p.id))
  return profiles
}

export type CreateStartFrom = Purpose | 'blank'

export function createStartOptions(): {
  id: CreateStartFrom
  title: string
  description: string
}[] {
  return [
    {
      id: 'general',
      title: 'General Assistant',
      description: 'Everyday questions and local-first chat.',
    },
    {
      id: 'coding',
      title: 'Programming',
      description: 'Code help with an AI team across your computers.',
    },
    {
      id: 'research',
      title: 'Research',
      description: 'Long-form reading and careful reasoning.',
    },
    {
      id: 'blank',
      title: 'Blank',
      description: 'Start empty and set roles, tools, and computers yourself.',
    },
  ]
}
