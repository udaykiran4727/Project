import { FormEvent, useEffect, useState } from 'react'
import { ApiReqError, createLink } from '../api/client'
import type { Link } from '../types/link'

interface LinkFormProps {
  onCreated: (link: Link) => void
  prefillShortcut?: string
}

export default function LinkForm({ onCreated, prefillShortcut }: LinkFormProps) {
  const [shortcut, setShortcut] = useState(prefillShortcut ?? '')
  const [destination, setDestination] = useState('')
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (prefillShortcut) {
      setShortcut(prefillShortcut)
    }
  }, [prefillShortcut])

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setFieldErrors({})
    setFormError(null)
    setSubmitting(true)

    try {
      const link = await createLink({ shortcut: shortcut.trim(), destination: destination.trim() })
      onCreated(link)
      setShortcut('')
      setDestination('')
    } catch (err) {
      if (err instanceof ApiReqError) {
        if (err.f) {
          setFieldErrors({ [err.f]: err.message })
        } else {
          setFormError(err.message)
        }
      } else {
        setFormError('Something went wrong. Is the backend running?')
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <h2 className="mb-4 text-lg font-semibold text-slate-800">New go link</h2>

      <div className="mb-4">
        <label htmlFor="shortcut" className="mb-1 block text-sm font-medium text-slate-700">
          Shortcut
        </label>
        <div className="flex items-center gap-2">
          <span className="text-sm text-slate-400">go/</span>
          <input
            id="shortcut"
            type="text"
            value={shortcut}
            onChange={(e) => setShortcut(e.target.value)}
            placeholder="oncall"
            aria-invalid={Boolean(fieldErrors.shortcut)}
            aria-describedby={fieldErrors.shortcut ? 'shortcut-error' : undefined}
            className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          />
        </div>
        {fieldErrors.shortcut && (
          <p id="shortcut-error" role="alert" className="mt-1 text-sm text-red-600">
            {fieldErrors.shortcut}
          </p>
        )}
      </div>

      <div className="mb-4">
        <label htmlFor="destination" className="mb-1 block text-sm font-medium text-slate-700">
          Destination URL
        </label>
        <input
          id="destination"
          type="text"
          value={destination}
          onChange={(e) => setDestination(e.target.value)}
          placeholder="https://example.com/oncall-schedule"
          aria-invalid={Boolean(fieldErrors.destination)}
          aria-describedby={fieldErrors.destination ? 'destination-error' : undefined}
          className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        {fieldErrors.destination && (
          <p id="destination-error" role="alert" className="mt-1 text-sm text-red-600">
            {fieldErrors.destination}
          </p>
        )}
      </div>

      {formError && (
        <div role="alert" className="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">
          {formError}
        </div>
      )}

      <button
        type="submit"
        disabled={submitting}
        className="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {submitting ? 'Creating…' : 'Create link'}
      </button>
    </form>
  )
}
