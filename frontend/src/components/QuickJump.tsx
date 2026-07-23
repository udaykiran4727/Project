import { FormEvent, useState } from 'react'
import { shortcutUrl } from '../api/client'

interface QuickJumpProps {
  onMiss: (shortcut: string) => void
}

/**
 * Lets a user type a shortcut and try to jump to it directly, mirroring
 * what typing "go/<shortcut>" in a browser address bar would do. If the
 * backend reports a 404, offers to prefill the create form instead.
 */
export default function QuickJump({ onMiss }: QuickJumpProps) {
  const [shortcut, setShortcut] = useState('')
  const [checking, setChecking] = useState(false)
  const [notFound, setNotFound] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    const trimmed = shortcut.trim()
    if (!trimmed) return

    setChecking(true)
    setNotFound(null)

    try {
      const res = await fetch(shortcutUrl(trimmed), { redirect: 'manual' })
      // An opaque response (status 0, type 'opaqueredirect') means the server
      // responded with a 3xx and the browser refused to expose it — that's
      // our signal the shortcut resolved. Navigate for real to follow it.
      if (res.type === 'opaqueredirect' || res.status === 0) {
        window.location.href = shortcutUrl(trimmed)
        return
      }
      if (res.status === 404) {
        setNotFound(trimmed)
        onMiss(trimmed)
        return
      }
      window.location.href = shortcutUrl(trimmed)
    } catch {
      setNotFound(trimmed)
    } finally {
      setChecking(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-2">
      <label htmlFor="quick-jump-shortcut" className="text-sm text-slate-400">
        go/
      </label>
      <input
        id="quick-jump-shortcut"
        type="text"
        value={shortcut}
        onChange={(e) => setShortcut(e.target.value)}
        placeholder="try a shortcut…"
        className="w-48 rounded-md border border-slate-300 px-3 py-1.5 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
      />
      <button
        type="submit"
        disabled={checking}
        className="rounded-md border border-slate-300 px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-100 disabled:opacity-50"
      >
        {checking ? 'Checking…' : 'Go'}
      </button>
      <span aria-live="polite" className="text-sm text-amber-700">
        {notFound && <>go/{notFound} doesn&apos;t exist yet — prefilled the form below.</>}
      </span>
    </form>
  )
}
