import { useEffect, useState } from 'react'
import LinkForm from './components/LinkForm'
import LinkTable from './components/LinkTable'
import QuickJump from './components/QuickJump'
import { listLinks } from './api/client'
import type { Link } from './types/link'

export default function App() {
  const [links, setLinks] = useState<Link[]>([])
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [prefillShortcut, setPrefillShortcut] = useState<string | undefined>(undefined)

  async function refresh() {
    setLoadError(null)
    try {
      const data = await listLinks()
      setLinks(data)
    } catch {
      setLoadError('Could not reach the go-links API. Is the backend running on port 8080?')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    refresh()
  }, [])

  function handleCreated(link: Link) {
    setLinks((prev) => [link, ...prev])
    setPrefillShortcut(undefined)
  }

  function handleDeleted(id: number) {
    setLinks((prev) => prev.filter((l) => l.id !== id))
  }

  return (
    <div className="min-h-screen">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-4xl items-center justify-between px-6 py-5">
          <div>
            <h1 className="text-xl font-bold text-slate-900">Go Links</h1>
            <p className="text-sm text-slate-500">Internal shortcut links for the team</p>
          </div>
          <QuickJump onMiss={setPrefillShortcut} />
        </div>
      </header>

      <main className="mx-auto max-w-4xl space-y-6 px-6 py-8">
        {loadError && (
          <div className="rounded-md bg-red-50 px-4 py-3 text-sm text-red-700">{loadError}</div>
        )}

        <LinkForm onCreated={handleCreated} prefillShortcut={prefillShortcut} />

        <div>
          <h2 className="mb-3 text-lg font-semibold text-slate-800">All links</h2>
          {loading ? (
            <p className="text-sm text-slate-500">Loading…</p>
          ) : (
            <LinkTable links={links} onDeleted={handleDeleted} />
          )}
        </div>
      </main>
    </div>
  )
}
