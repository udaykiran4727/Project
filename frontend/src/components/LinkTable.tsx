import { useState } from 'react'
import { deleteLink, shortcutUrl } from '../api/client'
import type { Link } from '../types/link'

interface LinkTableProps {
  links: Link[]
  onDeleted: (id: number) => void
}

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

export default function LinkTable({ links, onDeleted }: LinkTableProps) {
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)

  async function handleDelete(id: number) {
    setError(null)
    setDeletingId(id)
    try {
      await deleteLink(id)
      onDeleted(id)
    } catch {
      setError('Failed to delete link. Please try again.')
    } finally {
      setDeletingId(null)
    }
  }

  if (links.length === 0) {
    return (
      <div className="rounded-lg border border-dashed border-slate-300 bg-white p-8 text-center text-sm text-slate-500">
        No go links yet. Create one to get started.
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
      {error && (
        <div role="alert" className="bg-red-50 px-4 py-2 text-sm text-red-700">
          {error}
        </div>
      )}
      <table className="min-w-full divide-y divide-slate-200">
        <caption className="sr-only">Existing go links</caption>
        <thead className="bg-slate-50">
          <tr>
            <th scope="col" className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
              Shortcut
            </th>
            <th scope="col" className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
              Destination
            </th>
            <th scope="col" className="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide text-slate-500">
              Clicks
            </th>
            <th scope="col" className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
              Created
            </th>
            <th scope="col" className="px-4 py-3">
              <span className="sr-only">Actions</span>
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {links.map((link) => (
            <tr key={link.id} className="hover:bg-slate-50">
              <td className="px-4 py-3">
                <a
                  href={shortcutUrl(link.shortcut)}
                  target="_blank"
                  rel="noreferrer"
                  className="font-mono text-sm font-medium text-indigo-600 hover:underline"
                  title={`Follows the redirect to ${link.destination}`}
                >
                  go/{link.shortcut}
                </a>
              </td>
              <td className="max-w-xs truncate px-4 py-3 text-sm text-slate-600" title={link.destination}>
                {link.destination}
              </td>
              <td className="px-4 py-3 text-right text-sm text-slate-600">{link.click_count}</td>
              <td className="px-4 py-3 text-sm text-slate-500">{formatDate(link.created_at)}</td>
              <td className="px-4 py-3 text-right">
                <button
                  onClick={() => handleDelete(link.id)}
                  disabled={deletingId === link.id}
                  aria-label={`Delete go/${link.shortcut}`}
                  className="text-sm font-medium text-red-600 hover:text-red-800 disabled:opacity-50"
                >
                  {deletingId === link.id ? 'Deleting…' : 'Delete'}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
