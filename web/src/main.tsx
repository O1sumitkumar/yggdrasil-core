import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { App } from '@/app/App'
import '@/index.css'

async function bootstrap() {
  // Desktop (Wails) serves the UI and proxies /api to the daemon same-origin.
  // Leave __YGGDRASIL_API_BASE__ unset so fetch('/api/v1/...') stays same-origin
  // and avoids WebKit "Load failed" on POST (Private Network Access / CORS).
  // Vite dev mode uses the proxy in vite.config.ts instead.

  const rootElement = document.getElementById('root')
  if (!rootElement) {
    throw new Error('Root element #root not found')
  }

  createRoot(rootElement).render(
    <StrictMode>
      <App />
    </StrictMode>,
  )
}

void bootstrap()
