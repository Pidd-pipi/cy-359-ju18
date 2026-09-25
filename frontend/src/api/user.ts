import request from '../utils/request'

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    role: string
    points: number
  }
}

export interface User {
  id: number
  username: string
  nickname: string
  email?: string
  phone?: string
  role: string
  points: number
  disabled: boolean
  created_at: string
}

export function register(data: { username: string; password: string; nickname: string; email?: string; phone?: string }) {
  return request.post('/users/register', data)
}

export function login(data: { username: string; password: string }) {
  return request.post<LoginResponse>('/users/login', data)
}

export function getProfile() {
  return request.get<User>('/users/profile')
}

export function updateProfile(data: { nickname?: string; email?: string; phone?: string }) {
  return request.put('/users/profile', data)
}

export function listUsers(params: { page?: number; page_size?: number; keyword?: string }) {
  return request.get('/users', { params })
}

export function setUserDisabled(id: number, disabled: boolean) {
  return request.patch(`/users/${id}/disabled`, null, { params: { disabled } })
}
