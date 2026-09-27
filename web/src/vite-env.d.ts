/// <reference types="vite/client" />

interface YggdrasilWindow extends Window {
  __YGGDRASIL_API_BASE__?: string
  __YGGDRASIL_SCREENSHOT__?: { enabled?: boolean; screen?: string }
}
