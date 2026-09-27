import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ChatActivity } from './ChatActivity'

describe('ChatActivity', () => {
  it('announces that a reply is in progress', () => {
    render(<ChatActivity label="Searching the web…" />)
    expect(screen.getByRole('status', { name: 'Searching the web…' })).toBeInTheDocument()
    expect(screen.getByText('Searching the web…')).toBeInTheDocument()
  })
})
