/** Shown while a reply has not produced visible text yet. */
export function ChatActivity({ label }: { label?: string | null }) {
  return (
    <div
      className="flex max-w-[min(42rem,85%)] items-center gap-3 rounded-2xl bg-raised/80 px-4 py-3"
      role="status"
      aria-live="polite"
      aria-label={label || 'Working'}
    >
      <span className="inline-flex items-end gap-1" aria-hidden>
        <span className="chat-activity-dot" />
        <span className="chat-activity-dot" />
        <span className="chat-activity-dot" />
      </span>
      <span className="text-sm text-ink-muted">{label || 'Working…'}</span>
    </div>
  )
}
