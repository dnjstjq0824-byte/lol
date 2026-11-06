const TOKEN_KEY = 'riftstats_token'

export function getToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY)
  } catch {
    return null
  }
}

export function setToken(token: string | null): void {
  try {
    if (token) localStorage.setItem(TOKEN_KEY, token)
    else localStorage.removeItem(TOKEN_KEY)
  } catch {
    /* storage unavailable (private mode) */
  }
}

export interface ApiError {
  status: number
  title: string
  detail: string
}

export async function apiFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const headers: Record<string, string> = {
    ...(init.headers as Record<string, string> | undefined),
  }
  const token = getToken()
  if (token) headers['Authorization'] = `Bearer ${token}`
  return fetch(path, { ...init, headers })
}

export async function apiJson<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await apiFetch(path, init)
  if (!res.ok) {
    let title = res.statusText
    let detail = `Request failed with status ${res.status}`
    try {
      const problem = await res.json()
      if (problem?.title) title = problem.title
      if (problem?.detail) detail = problem.detail
    } catch {
      /* non-JSON error body */
    }
    const err = new Error(detail) as Error & ApiError
    err.status = res.status
    err.title = title
    err.detail = detail
    throw err
  }
  return res.json() as Promise<T>
}

export function isApiError(e: unknown): e is Error & ApiError {
  return !!e && typeof e === 'object' && 'status' in (e as Record<string, unknown>)
}
