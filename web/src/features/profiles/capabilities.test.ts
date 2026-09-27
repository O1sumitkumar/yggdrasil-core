import { describe, expect, it } from 'vitest'
import { capabilityEnabled, setCapability } from './capabilities'
import type { ToolPolicy } from '@/types/api'

const base: ToolPolicy[] = [
  { tool_id: 'internet.search', policy: 'allow' },
  { tool_id: 'internet.open', policy: 'allow' },
  { tool_id: 'terminal', policy: 'deny' },
]

describe('capabilities', () => {
  it('treats internet as on only when a web tool is allowed', () => {
    expect(capabilityEnabled(base, 'internet')).toBe(true)
    expect(capabilityEnabled(base, 'shell')).toBe(false)
  })

  it('switches a capability without inventing shell access', () => {
    const off = setCapability(base, 'internet', false)
    expect(off.find((tool) => tool.tool_id === 'internet.search')?.policy).toBe('deny')
    expect(off.find((tool) => tool.tool_id === 'terminal')?.policy).toBe('deny')
    const shell = setCapability(off, 'shell', true)
    expect(shell.find((tool) => tool.tool_id === 'terminal')?.policy).toBe('allow')
  })
})