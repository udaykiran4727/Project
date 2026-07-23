export interface Link {
  id: number
  shortcut: string
  destination: string
  created_at: string
  click_count: number
}

export interface CreateLinkRequest {
  shortcut: string
  destination: string
}

export interface ApiError {
  error: string
  field?: string
}
