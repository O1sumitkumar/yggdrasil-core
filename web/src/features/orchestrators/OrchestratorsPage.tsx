import { EmptyState } from '@/components/ui/EmptyState'

export function OrchestratorsPage() {
  return (
    <div className="w-full min-w-0 space-y-6">
      <header>
        <h1 className="font-display text-3xl font-semibold text-ink">Orchestrators</h1>
        <p className="mt-1 text-sm text-ink-muted">
          Pluggable orchestration engines that coordinate models and tools.
        </p>
      </header>

      <EmptyState
        title="Orchestrator management coming soon"
        description="Built-in and plugin orchestrators will be configurable here under advanced mode."
      />
    </div>
  )
}
