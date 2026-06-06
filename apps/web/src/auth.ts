export type ApiEnvelope<T> = {
  code: number
  msg: string
  data: T
}

export type AuthUser = {
  id: string
  email: string
  display_name: string
}

type AuthData = {
  user: AuthUser
}

type SendCodeData = {
  cooldown_seconds: number
  expires_in_seconds: number
}

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') ?? ''

async function request<T>(path: string, options: RequestInit = {}) {
  const headers = new Headers(options.headers)

  if (!headers.has('Content-Type') && options.body) {
    headers.set('Content-Type', 'application/json')
  }

  const response = await fetch(`${apiBaseUrl}${path}`, { ...options, headers, credentials: 'include' })
  const result = (await response.json()) as ApiEnvelope<T>

  if (!response.ok || result.code !== 200) {
    throw new Error(result.msg || 'Request failed')
  }

  return result.data
}

export const sendRegisterCode = async (payload: { email: string }) => {
  return request<SendCodeData>('/api/v1/auth/send-code', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const login = async (payload: { email: string; password: string }) => {
  return request<AuthData>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const register = async (payload: { name: string; email: string; password: string; verification_code: string }) => {
  return request<AuthData>('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const logout = () => {
  return request<Record<string, never>>('/api/v1/auth/logout', { method: 'POST' })
}

export const fetchMe = () => request<AuthUser>('/api/v1/auth/me')
