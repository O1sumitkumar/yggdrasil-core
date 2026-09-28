import { useCallback, useEffect, useRef, useState } from 'react'

/** Keep following while the reader is within this distance of the newest line. */
export const CHAT_BOTTOM_THRESHOLD_PX = 72

export function distanceFromBottom(el: {
  scrollHeight: number
  scrollTop: number
  clientHeight: number
}): number {
  return el.scrollHeight - el.scrollTop - el.clientHeight
}

export function shouldFollowScroll(
  distance: number,
  threshold = CHAT_BOTTOM_THRESHOLD_PX,
): boolean {
  return distance <= threshold
}

/**
 * Keep a transcript pinned to the newest content until the reader scrolls up.
 * A pinned view follows growth; an unpinned view stays put and can jump back.
 */
export function useChatFollow(conversationId: string | null) {
  const scrollerRef = useRef<HTMLDivElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const stickRef = useRef(true)
  const ignoreScrollRef = useRef(false)
  const [showJump, setShowJump] = useState(false)

  const applyStick = useCallback((stick: boolean) => {
    if (stickRef.current === stick) return
    stickRef.current = stick
    setShowJump(!stick)
  }, [])

  const scrollToEnd = useCallback((behavior: ScrollBehavior) => {
    const el = scrollerRef.current
    if (!el) return
    const top = el.scrollHeight
    if (shouldFollowScroll(distanceFromBottom(el))) {
      ignoreScrollRef.current = false
      return
    }
    ignoreScrollRef.current = true
    if (behavior === 'auto') {
      el.scrollTop = top
      ignoreScrollRef.current = false
      return
    }
    el.scrollTo({ top, behavior })
  }, [])

  const followLatest = useCallback(() => {
    applyStick(true)
    scrollToEnd('auto')
  }, [applyStick, scrollToEnd])

  const jumpToLatest = useCallback(() => {
    applyStick(true)
    scrollToEnd('smooth')
  }, [applyStick, scrollToEnd])

  useEffect(() => {
    stickRef.current = true
    setShowJump(false)
    if (!conversationId) return

    const scroller = scrollerRef.current
    const content = contentRef.current
    if (!scroller || !content) return

    const followIfStuck = () => {
      if (!stickRef.current) return
      const el = scrollerRef.current
      if (!el || shouldFollowScroll(distanceFromBottom(el), 1)) return
      ignoreScrollRef.current = true
      el.scrollTop = el.scrollHeight
      ignoreScrollRef.current = false
    }

    const onScroll = () => {
      if (ignoreScrollRef.current) {
        if (shouldFollowScroll(distanceFromBottom(scroller))) {
          ignoreScrollRef.current = false
          applyStick(true)
        }
        return
      }
      applyStick(shouldFollowScroll(distanceFromBottom(scroller)))
    }

    const releaseProgrammatic = () => {
      const el = scrollerRef.current
      const wasAnimating = ignoreScrollRef.current
      ignoreScrollRef.current = false
      if (wasAnimating && el) {
        // Reassigning scrollTop cancels an in-progress smooth scroll.
        // eslint-disable-next-line no-self-assign
        el.scrollTop = el.scrollTop
      }
      if (el) applyStick(shouldFollowScroll(distanceFromBottom(el)))
    }

    const onWheel = (event: WheelEvent) => {
      if (event.deltaY >= 0) return
      if (scroller.scrollHeight <= scroller.clientHeight + 1) return
      ignoreScrollRef.current = false
      applyStick(false)
    }

    const observer = new ResizeObserver(() => {
      followIfStuck()
    })
    observer.observe(content)
    scroller.addEventListener('scroll', onScroll, { passive: true })
    scroller.addEventListener('wheel', onWheel, { passive: true })
    scroller.addEventListener('pointerdown', releaseProgrammatic)
    followIfStuck()

    return () => {
      observer.disconnect()
      scroller.removeEventListener('scroll', onScroll)
      scroller.removeEventListener('wheel', onWheel)
      scroller.removeEventListener('pointerdown', releaseProgrammatic)
    }
  }, [applyStick, conversationId])

  return { scrollerRef, contentRef, showJump, jumpToLatest, followLatest }
}
