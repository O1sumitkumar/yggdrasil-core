export function TightFitDialog({
  modelName,
  onCancel,
  onInstall,
}: {
  modelName: string
  onCancel: () => void
  onInstall: () => void
}) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center bg-ink/40 p-4 sm:items-center"
      role="presentation"
      onClick={onCancel}
    >
      <div
        className="card w-full max-w-md space-y-4 border-l-4 border-danger shadow-panel"
        role="dialog"
        aria-modal="true"
        aria-labelledby="tight-fit-title"
        onClick={(event) => event.stopPropagation()}
      >
        <div className="flex gap-3">
          <span
            className="mt-0.5 inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-danger/15 text-sm font-semibold text-danger"
            aria-hidden
          >
            !
          </span>
          <div>
            <h2 id="tight-fit-title" className="font-display text-lg font-semibold text-ink">
              This model may be unstable on this computer
            </h2>
            <p className="mt-2 text-sm text-ink-muted">
              {modelName} uses most of the memory available on this computer. It may run slowly,
              stop responding, or cause system instability under heavier workloads or larger context
              sizes.
            </p>
            <p className="mt-2 text-sm text-ink-muted">
              Yggdrasil recommends leaving more memory available for the operating system and graphics.
            </p>
          </div>
        </div>
        <div className="flex flex-wrap justify-end gap-2">
          <button type="button" className="btn-secondary" autoFocus onClick={onCancel}>
            Cancel
          </button>
          <button type="button" className="btn-primary" onClick={onInstall}>
            Install anyway
          </button>
        </div>
      </div>
    </div>
  )
}
