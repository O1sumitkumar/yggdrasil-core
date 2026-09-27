import { act, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { distanceFromBottom, shouldFollowScroll, useChatFollow } from './useChatFollow'

class FakeResizeObserver {
  static instances: FakeResizeObserver[] = []
  private readonly callback: ResizeObserverCallback

  constructor(callback: ResizeObserverCallback) {
    this.callback = callback
    FakeResizeObserver.instances.push(this)
  }

  observe() {}
  unobserve() {}
  disconnect() {}

  grow() {
    this.callback([], this)
  }
}

type Metrics = {
  scrollHeight: number
  clientHeight: number
  scrollTop: number
}

function installMetrics(el: HTMLElement, metrics: Metrics) {
  Object.defineProperty(el, 'scrollHeight', {
    configurable: true,
    get: () => metrics.scrollHeight,
  })
  Object.defineProperty(el, 'clientHeight', {
    configurable: true,
    get: () => metrics.clientHeight,
  })
  Object.defineProperty(el, 'scrollTop', {
    configurable: true,
    get: () => metrics.scrollTop,
    set: (value: number) => {
      const max = Math.max(0, metrics.scrollHeight - metrics.clientHeight)
      metrics.scrollTop = Math.max(0, Math.min(value, max))
      el.dispatchEvent(new Event('scroll'))
    },
  })
  el.scrollTo = (optionsOrX?: ScrollToOptions | number, y?: number) => {
    const top = typeof optionsOrX === 'number' ? y : optionsOrX?.top
    if (typeof top === 'number') el.scrollTop = top
  }
}

function Transcript({ conversationId }: { conversationId: string | null }) {
  const { scrollerRef, contentRef, showJump, jumpToLatest } = useChatFollow(conversationId)
  return (
    <div>
      <div ref={scrollerRef} data-testid="scroller">
        <div ref={contentRef} />
      </div>
      {showJump ? (
        <button type="button" onClick={jumpToLatest}>
          Latest
        </button>
      ) : null}
    </div>
  )
}

describe('shouldFollowScroll', () => {
  it('stays with the newest lines inside the threshold', () => {
    expect(distanceFromBottom({ scrollHeight: 500, scrollTop: 300, clientHeight: 180 })).toBe(20)
    expect(shouldFollowScroll(20)).toBe(true)
    expect(shouldFollowScroll(200)).toBe(false)
  })
})

describe('useChatFollow', () => {
  it('follows new content until the reader scrolls up, then offers a jump', () => {
    vi.stubGlobal('ResizeObserver', FakeResizeObserver)
    FakeResizeObserver.instances = []

    const metrics: Metrics = { scrollHeight: 400, clientHeight: 200, scrollTop: 200 }
    const view = render(<Transcript conversationId="chat-1" />)
    const scroller = screen.getByTestId('scroller')
    installMetrics(scroller, metrics)

    act(() => {
      FakeResizeObserver.instances.at(-1)?.grow()
    })
    expect(metrics.scrollTop).toBe(200)
    expect(screen.queryByRole('button', { name: 'Latest' })).toBeNull()

    metrics.scrollHeight = 900
    act(() => {
      FakeResizeObserver.instances.at(-1)?.grow()
    })
    expect(metrics.scrollTop).toBe(700)
    expect(screen.queryByRole('button', { name: 'Latest' })).toBeNull()

    act(() => {
      scroller.scrollTop = 100
    })
    expect(screen.getByRole('button', { name: 'Latest' })).toBeTruthy()

    metrics.scrollHeight = 1200
    act(() => {
      FakeResizeObserver.instances.at(-1)?.grow()
    })
    expect(metrics.scrollTop).toBe(100)

    act(() => {
      screen.getByRole('button', { name: 'Latest' }).click()
    })
    expect(metrics.scrollTop).toBe(1000)
    expect(screen.queryByRole('button', { name: 'Latest' })).toBeNull()

    view.unmount()
    vi.unstubAllGlobals()
  })
})
