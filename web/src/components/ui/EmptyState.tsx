import type { ReactNode } from 'react'

interface EmptyStateProps {
  title: string
  description: string
  action?: ReactNode
}

export function EmptyState({ title, description, action }: EmptyStateProps) {
  return (
    <div className="empty-hero max-w-lg">
      <h3 className="font-display text-xl font-semibold text-ink">{title}</h3>
      <p className="text-[15px] leading-relaxed text-ink-muted">{description}</p>
      {action ? <div className="pt-2">{action}</div> : null}
    </div>
  )
}
