import { Link } from 'react-router-dom'
import type { CapabilityGap } from '@/features/models/capabilityGap'

export function CapabilityNotice({
  gap,
  onUse,
  onEnableInternet,
}: {
  gap: CapabilityGap
  onUse: (modelId: string) => void
  onEnableInternet?: () => void
}) {
  return (
    <div className="max-w-[min(42rem,85%)] space-y-3 rounded-2xl border border-line/70 bg-raised/60 px-4 py-3 text-sm text-ink">
      {gap.notes.map((note) => (
        <div key={note.text}>
          <p>{note.text}</p>
          {note.action === 'enable-internet' && onEnableInternet && (
            <button
              type="button"
              className="mt-2 text-xs font-medium text-primary underline-offset-2 hover:underline"
              onClick={onEnableInternet}
            >
              Enable Internet
            </button>
          )}
          {note.suggestions.length > 0 && (
            <ul className="mt-2 space-y-1.5">
              {note.suggestions.map((item) => (
                <li key={item.id} className="flex flex-wrap items-center gap-x-2 gap-y-1">
                  <span className="font-medium">{item.name}</span>
                  <span className="text-ink-faint">{item.installed ? 'Installed' : 'Not installed'}</span>
                  {item.installed ? (
                    <button
                      type="button"
                      className="text-xs font-medium text-primary underline-offset-2 hover:underline"
                      onClick={() => onUse(item.id)}
                    >
                      Use this model
                    </button>
                  ) : (
                    <Link
                      to="/models"
                      className="text-xs font-medium text-primary underline-offset-2 hover:underline"
                    >
                      View in Models
                    </Link>
                  )}
                </li>
              ))}
            </ul>
          )}
        </div>
      ))}
    </div>
  )
}
