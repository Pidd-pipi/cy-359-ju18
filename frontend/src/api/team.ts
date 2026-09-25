import request from '../utils/request'

export interface TeamMember {
  user_id: number
  username: string
  nickname: string
  role: string
}

export interface Team {
  id: number
  name: string
  slogan?: string
  captain_id: number
  member_count: number
  members?: TeamMember[]
  created_at: string
}

export interface Registration {
  id: number
  team_id: number
  team_name: string
  activity_id: number
  activity_title: string
  status: string
  start_time?: string
  finish_time?: string
  total_seconds: number
  duration: string
  registered_at: string
}

export function createTeam(data: { name: string; slogan?: string }) {
  return request.post('/teams', data)
}

export function listMyTeams() {
  return request.get<Team[]>('/teams/mine')
}

export function listTeams(params: { page?: number; page_size?: number }) {
  return request.get<{ list: Team[]; total: number }>('/teams', { params })
}

export function getTeam(id: number) {
  return request.get<Team>(`/teams/${id}`)
}

export function joinTeam(teamId: number) {
  return request.post(`/teams/${teamId}/join`, { team_id: teamId })
}

export function leaveTeam(teamId: number) {
  return request.delete(`/teams/${teamId}/leave`)
}

export function applyActivity(data: { team_id: number; activity_id: number }) {
  return request.post('/registrations', data)
}

export function listMyRegistrations() {
  return request.get<Registration[]>('/registrations/mine')
}

export function listRegistrationsByActivity(activityId: number) {
  return request.get<Registration[]>(`/activities/${activityId}/registrations`)
}

export function approveRegistration(id: number) {
  return request.post(`/registrations/${id}/approve`)
}

export function rejectRegistration(id: number) {
  return request.post(`/registrations/${id}/reject`)
}
