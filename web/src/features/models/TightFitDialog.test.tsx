import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { TightFitDialog } from './TightFitDialog'

describe('TightFitDialog', () => {
  it('does not install when cancelled', () => {
    const onInstall = vi.fn()
    const onCancel = vi.fn()
    render(<TightFitDialog modelName="Qwen 32B" onCancel={onCancel} onInstall={onInstall} />)
    expect(screen.getByRole('dialog')).toHaveTextContent('may be unstable')
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(onCancel).toHaveBeenCalledOnce()
    expect(onInstall).not.toHaveBeenCalled()
  })

  it('installs when the override is confirmed', () => {
    const onInstall = vi.fn()
    render(<TightFitDialog modelName="Qwen 32B" onCancel={() => {}} onInstall={onInstall} />)
    fireEvent.click(screen.getByRole('button', { name: 'Install anyway' }))
    expect(onInstall).toHaveBeenCalledOnce()
  })
})
