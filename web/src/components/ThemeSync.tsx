import { useEffect } from 'react'
import { applyTheme, useUIStore } from '@/stores/uiStore'

/** Keeps data-theme in sync with preference (including system changes). */
export function ThemeSync() {
  const theme = useUIStore((s) => s.theme)

  useEffect(() => {
    applyTheme(theme)
    if (theme !== 'system') {
      return
    }
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => applyTheme('system')
    mq.addEventListener('change', onChange)
    return () => mq.removeEventListener('change', onChange)
  }, [theme])

  return null
}
