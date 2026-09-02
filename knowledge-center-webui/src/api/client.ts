export type ImportStatus = 'pending' | 'running' | 'success' | 'failed'

export interface ImportItem {
  id: string
  url: string
  host: string
  status: ImportStatus
  error_message: string
  created_at: string
  updated_at: string
  scraped_at?: string | null
}

export interface PageSummary {
  id: string
  import_id: string
  url: string
  title: string
  created_at: string
  scraped_at: string
}

export interface ImportDetail extends ImportItem {
  pages: PageSummary[]
}

export interface PageDetail extends PageSummary {
  content_text: string
  content_html: string
}

export interface ScrapeProfile {
  id: string
  host: string
  title_selector: string
  content_selector: string
  link_selector: string
  exclude_selector: string
  created_at: string
  updated_at: string
}

export type ProfileInput = {
  host: string
  title_selector: string
  content_selector: string
  link_selector: string
  exclude_selector: string
}

const baseUrl = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') ?? ''

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${baseUrl}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
    ...init,
  })

  if (response.status === 204) {
    return undefined as T
  }

  const text = await response.text()
  let data: unknown = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      throw new Error(text.trim() || `Request failed (${response.status})`)
    }
  }

  if (!response.ok) {
    const message =
      typeof data === 'object' && data && 'error' in data
        ? String((data as { error: string }).error)
        : `Request failed (${response.status})`
    throw new Error(message)
  }

  return data as T
}

export function listImports() {
  return request<ImportItem[]>('/api/imports')
}

export function createImport(url: string) {
  return request<ImportItem>('/api/imports', {
    method: 'POST',
    body: JSON.stringify({ url }),
  })
}

export function getImport(id: string) {
  return request<ImportDetail>(`/api/imports/${id}`)
}

export function getPage(importId: string, pageId: string) {
  return request<PageDetail>(`/api/imports/${importId}/pages/${pageId}`)
}

export function listProfiles() {
  return request<ScrapeProfile[]>('/api/scrape-profiles')
}

export function createProfile(input: ProfileInput) {
  return request<ScrapeProfile>('/api/scrape-profiles', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function updateProfile(id: string, input: ProfileInput) {
  return request<ScrapeProfile>(`/api/scrape-profiles/${id}`, {
    method: 'PUT',
    body: JSON.stringify(input),
  })
}

export function deleteProfile(id: string) {
  return request<void>(`/api/scrape-profiles/${id}`, { method: 'DELETE' })
}

export function health() {
  return request<{ status: string }>('/api/health')
}
