import { Outlet } from 'react-router-dom'
import { Sidebar } from '@/components/layout/Sidebar'

/**
 * App chrome: sidebar + main page region.
 * Overflow is owned here — pages scroll inside `.page-scroll`, not the OS window.
 */
export function AppLayout() {
  return (
    <div className="app-shell">
      <Sidebar />
      <main className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden">
        <div className="page-scroll px-page-x py-5 sm:py-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
