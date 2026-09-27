import { spawn } from 'node:child_process'
import { mkdir, readFile } from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright'

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..')
const config = JSON.parse(await readFile(path.join(root, 'screenshots/config.json'), 'utf8'))
const width = config.window?.width || 1440
const height = config.window?.height || 900
const outDir = path.join(root, 'docs/screenshots')
await mkdir(outDir, { recursive: true })

const server = spawn('go', ['run', './cmd/screenshot', '-web', 'web/dist'], {
  cwd: root,
  detached: true,
  stdio: ['ignore', 'pipe', 'inherit'],
})
let base = ''
let log = ''
const ready = new Promise((resolve, reject) => {
  const timer = setTimeout(() => reject(new Error('screenshot server did not start')), 120_000)
  server.stdout.on('data', (chunk) => {
    log += chunk.toString()
    const match = log.match(/listening (http:\/\/\S+)/)
    if (match) {
      clearTimeout(timer)
      resolve(match[1])
    }
  })
  server.on('exit', (code) => {
    if (!base) reject(new Error(`screenshot server exited ${code}`))
  })
})

try {
  base = await ready
  const browser = await chromium.launch(
    process.env.SCREENSHOT_CHROME ? { executablePath: process.env.SCREENSHOT_CHROME } : {},
  )
  const page = await browser.newPage({ viewport: { width, height }, deviceScaleFactor: 1 })
  for (const screen of config.screens) {
    const query = new URLSearchParams({ screenshot: '1', screen: screen.id })
    if (screen.query) {
      for (const [key, value] of new URLSearchParams(screen.query)) {
        query.set(key, value)
      }
    }
    const url = `${base}${screen.path}?${query}`
    console.log(`Capturing ${screen.id} ${url}`)
    await page.goto(url, { waitUntil: 'domcontentloaded' })
    try {
      await page.waitForFunction(
        (text) => document.body?.innerText?.includes(text),
        screen.readyText,
        { timeout: 20_000 },
      )
    } catch (error) {
      const text = await page.locator('body').innerText().catch(() => '')
      console.error(`Page text for ${screen.id}:\n${text.slice(0, 1200)}`)
      throw error
    }
    await page.screenshot({ path: path.join(outDir, screen.filename) })
  }
  await browser.close()
  console.log(`Screenshots: ${outDir}`)
} finally {
  try {
    process.kill(-server.pid, 'SIGTERM')
  } catch {
    server.kill()
  }
}
