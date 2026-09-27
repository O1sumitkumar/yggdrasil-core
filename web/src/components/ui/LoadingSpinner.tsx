export function LoadingSpinner({ label = 'Loading…' }: { label?: string }) {
  return (
    <div className="flex items-center gap-3 text-sm text-ink-muted" role="status">
      <span className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-line border-t-primary" />
      <span>{label}</span>
    </div>
  )
}
