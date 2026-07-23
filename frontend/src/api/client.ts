import type { ApiError, CreateLinkRequest, Link } from '../types/link'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export class ApiReqError extends Error {
  s: number
  f?: string

  constructor(s: number, msg: string, f?: string) {
    super(msg)
    this.s = s
    this.f = f
  }
}

async function handleResponse<T>(res: Response): Promise<T> {
  if (res.ok) {
    if (res.status === 204) {
      return undefined as T
    }
    return (await res.json()) as T
  }

  let body: ApiError | undefined
  try {
    body = (await res.json()) as ApiError
  } catch {
    // response had no JSON body
  }
  throw new ApiReqError(res.status, body?.error ?? res.statusText, body?.field)
}

export async function listLinks(): Promise<Link[]> {
  const res = await fetch(`${API_BASE}/api/links`)
  return handleResponse<Link[]>(res)
}

export async function getLink(id: number): Promise<Link> {
  const res = await fetch(`${API_BASE}/api/links/${id}`)
  return handleResponse<Link>(res)
}

export async function createLink(req: CreateLinkRequest): Promise<Link> {
  const res = await fetch(`${API_BASE}/api/links`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  return handleResponse<Link>(res)
}

export async function deleteLink(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/api/links/${id}`, { method: 'DELETE' })
  return handleResponse<void>(res)
}

export function shortcutUrl(shortcut: string): string {
  return `${API_BASE}/${shortcut}`
}
