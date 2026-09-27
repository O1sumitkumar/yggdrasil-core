/** Optional bridge to the Wails desktop shell (no-op in the browser). */

type WailsApp = {
  SetBackgroundMode?: (v: boolean) => Promise<void>
  QuitAndStopDaemon?: () => Promise<void>
  OpenPath?: (path: string) => Promise<void>
  SetLaunchAtLogin?: (v: boolean) => Promise<void>
  IsLaunchAtLogin?: () => Promise<boolean>
  MarkScreenshotReady?: () => Promise<void>
}

type WailsRuntime = {
  EventsOn?: (eventName: string, callback: (...data: unknown[]) => void) => () => void
}

function wailsGoApp(): WailsApp | undefined {
  return (window as unknown as { go?: { main?: { App?: WailsApp } } }).go?.main?.App
}

function wailsRuntime(): WailsRuntime | undefined {
  return (window as unknown as { runtime?: WailsRuntime }).runtime
}

export async function notifyDesktopBackgroundMode(enabled: boolean): Promise<void> {
  try {
    await wailsGoApp()?.SetBackgroundMode?.(enabled)
  } catch {
    // Running in a normal browser against the daemon — ignore.
  }
}

/** Enable or disable OS launch-at-login when running in the desktop shell. */
export async function notifyDesktopLaunchAtLogin(enabled: boolean): Promise<void> {
  try {
    await wailsGoApp()?.SetLaunchAtLogin?.(enabled)
  } catch {
    // Browser / headless — setting is still stored by the daemon.
  }
}

export function isDesktopShell(): boolean {
  return Boolean(wailsGoApp())
}

export async function markScreenshotReady(): Promise<void> {
  try {
    await wailsGoApp()?.MarkScreenshotReady?.()
  } catch {
    // The capture script also watches data-screenshot-ready on the page.
  }
}

/** Reveal a path in Finder / Explorer / file manager when running in desktop. */
export async function openPathInOS(path: string): Promise<boolean> {
  if (!path) return false
  try {
    const app = wailsGoApp()
    if (app?.OpenPath) {
      await app.OpenPath(path)
      return true
    }
  } catch {
    // fall through
  }
  try {
    const runtime = wailsRuntime() as WailsRuntime & {
      BrowserOpenURL?: (url: string) => void
    }
    if (runtime?.BrowserOpenURL) {
      const href = path.startsWith('file:') ? path : `file://${path}`
      runtime.BrowserOpenURL(href)
      return true
    }
  } catch {
    // ignore
  }
  return false
}

/** Ask the desktop shell to quit and stop the daemon (full restart by relaunching). */
export async function quitDesktopForRestart(): Promise<boolean> {
  try {
    const app = wailsGoApp()
    if (app?.QuitAndStopDaemon) {
      await app.QuitAndStopDaemon()
      return true
    }
  } catch {
    // ignore
  }
  return false
}

export type DesktopLifecycleEvent = {
  state: 'stopping' | 'background' | string
  message?: string
}

/** Subscribe to desktop shell lifecycle events. Returns an unsubscribe fn. */
export function onDesktopLifecycle(
  handler: (event: DesktopLifecycleEvent) => void,
): () => void {
  const runtime = wailsRuntime()
  if (!runtime?.EventsOn) {
    return () => {}
  }
  return runtime.EventsOn('ygg:lifecycle', (...data: unknown[]) => {
    const raw = data[0]
    if (raw && typeof raw === 'object') {
      const obj = raw as Record<string, unknown>
      handler({
        state: typeof obj.state === 'string' ? obj.state : 'stopping',
        message: typeof obj.message === 'string' ? obj.message : undefined,
      })
      return
    }
    handler({ state: 'stopping' })
  })
}
