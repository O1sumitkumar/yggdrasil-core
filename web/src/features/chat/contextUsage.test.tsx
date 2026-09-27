import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ContextUsageButton } from './ContextUsageButton'
import { fillPercent, formatTokens, parseContextUsage } from './contextUsage'

describe('context usage formatting', () => {
  it('formats token counts and percent full', () => {
    expect(formatTokens(2100)).toBe('2.1K')
    expect(formatTokens(8192)).toBe('8.2K')
    expect(formatTokens(256000)).toBe('256K')
    expect(fillPercent(95800, 256000)).toBe(37)
  })

  it('reads the server breakdown', () => {
    const usage = parseContextUsage({
      prompt_tokens: 100,
      limit: 8192,
      instructions: 10,
      tools: 20,
      conversation: 30,
      tool_results: 40,
      estimated: false,
    })
    expect(usage).toMatchObject({
      promptTokens: 100,
      tools: 20,
      toolResults: 40,
    })
  })
})

describe('ContextUsageButton', () => {
  it('opens the breakdown of what the model is holding', () => {
    render(
      <ContextUsageButton
        windowLimit={8192}
        usage={{
          promptTokens: 100,
          limit: 8192,
          instructions: 10,
          tools: 25,
          conversation: 40,
          toolResults: 25,
          estimated: false,
        }}
      />,
    )
    fireEvent.click(screen.getByRole('button', { name: 'Context 1% full' }))
    expect(screen.getByText('1% full')).toBeTruthy()
    expect(screen.getByText('Instructions')).toBeTruthy()
    expect(screen.getByText('Tools')).toBeTruthy()
    expect(screen.getByText('Conversation')).toBeTruthy()
    expect(screen.getByText('Tool results')).toBeTruthy()
    expect(screen.getByText('100 / 8.2K tokens')).toBeTruthy()
  })
})
