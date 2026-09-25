// 本地登录态读写
export interface AuthUser {
  id: number
  username: string
  nickname: string
  role: string
  points: number
}

export function getToken(): string | null {
  return localStorage.getItem('orienteering_token')
}

export function setAuth(token: string, user: AuthUser) {
  localStorage.setItem('orienteering_token', token)
  localStorage.setItem('orienteering_user', JSON.stringify(user))
}

export function getUser(): AuthUser | null {
  const raw = localStorage.getItem('orienteering_user')
  if (!raw) return null
  try {
    return JSON.parse(raw) as AuthUser
  } catch {
    return null
  }
}

export function clearAuth() {
  localStorage.removeItem('orienteering_token')
  localStorage.removeItem('orienteering_user')
}

export function isAdmin(): boolean {
  return getUser()?.role === 'admin'
}
