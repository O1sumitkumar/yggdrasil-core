import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { markScreenshotReady } from '@/lib/desktopBridge'
import { readScreenshotLaunch } from '@/lib/screenshotMode'
import { useUIStore } from '@/stores/uiStore'

/** Opens the requested screen and signals the capture script after paint. */
export function ScreenshotMode() {
  const navigate = useNavigate()

  useEffect(() => {
    const launch = readScreenshotLaunch()
    if (!launch?.enabled) {
      return
    }
    document.documentElement.dataset.screenshot = '1'
    useUIStore.setState({ onboardingComplete: true, advancedMode: true, theme: 'dark' })
    useUIStore.getState().setTheme('dark')
    if (window.location.pathname !== launch.path) {
      navigate(launch.path, { replace: true })
    }
    const stopHydration = useUIStore.persist.onFinishHydration(() => {
      useUIStore.setState({ onboardingComplete: true, advancedMode: true, theme: 'dark' })
    })
    const timer = window.setTimeout(() => {
      void markScreenshotReady()
      document.documentElement.dataset.screenshotReady = '1'
    }, 1200)
    return () => {
      stopHydration()
      window.clearTimeout(timer)
    }
  }, [navigate])

  return null
}
