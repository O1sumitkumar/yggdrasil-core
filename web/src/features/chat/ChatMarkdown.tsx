import type { Components } from 'react-markdown'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { prepareChatMarkdown } from './displayChatText'

const components: Components = {
  a: ({ href, children }) => {
    if (!href || !/^https?:\/\//i.test(href)) return <span>{children}</span>
    return (
      <a href={href} target="_blank" rel="noopener noreferrer" className="text-primary underline-offset-2 hover:underline">
        {children}
      </a>
    )
  },
  p: ({ children }) => <p className="mb-2 last:mb-0">{children}</p>,
  ul: ({ children }) => <ul className="my-2 list-disc space-y-1 pl-5 last:mb-0">{children}</ul>,
  ol: ({ children }) => <ol className="my-2 list-decimal space-y-1 pl-5 last:mb-0">{children}</ol>,
  li: ({ children }) => <li className="pl-0.5">{children}</li>,
  h1: ({ children }) => <h2 className="mb-2 font-display text-lg font-semibold">{children}</h2>,
  h2: ({ children }) => <h3 className="mb-2 font-display text-base font-semibold">{children}</h3>,
  h3: ({ children }) => <h4 className="mb-1 font-semibold">{children}</h4>,
  blockquote: ({ children }) => (
    <blockquote className="my-2 border-l-2 border-primary/40 pl-3 text-ink-muted">{children}</blockquote>
  ),
  pre: ({ children }) => (
    <pre className="my-2 overflow-x-auto rounded-lg bg-canvas p-3 font-mono text-xs leading-relaxed last:mb-0">
      {children}
    </pre>
  ),
  code: ({ className, children }) => {
    if (className) return <code className={className}>{children}</code>
    return <code className="rounded bg-canvas/80 px-1 py-0.5 font-mono text-[0.9em]">{children}</code>
  },
  table: ({ children }) => (
    <div className="my-2 overflow-x-auto">
      <table className="w-full border-collapse text-sm">{children}</table>
    </div>
  ),
  th: ({ children }) => <th className="border-b border-line px-2 py-1 text-left font-medium">{children}</th>,
  td: ({ children }) => <td className="border-b border-line/60 px-2 py-1 align-top">{children}</td>,
}

export function ChatMarkdown({ text }: { text: string }) {
  return (
    <div className="chat-markdown">
      <ReactMarkdown remarkPlugins={[remarkGfm]} components={components}>
        {prepareChatMarkdown(text)}
      </ReactMarkdown>
    </div>
  )
}
