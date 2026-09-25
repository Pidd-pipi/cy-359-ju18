import request from '../utils/request'

export interface Checkpoint {
  id: number
  activity_id: number
  name: string
  sequence: number
  lat: number
  lng: number
  clue?: string
  task_type: string
  task_content?: string
  radius_meters: number
  qr_code: string
}

export interface Activity {
  id: number
  title: string
  description?: string
  difficulty: string
  duration_minutes: number
  equipment_requirement?: string
  start_time: string
  end_time: string
  status: string
  creator_id: number
  start_lat: number
  start_lng: number
  end_lat: number
  end_lng: number
  address: string
  max_teams: number
  checkpoint_count: number
  team_count: number
  checkpoints?: Checkpoint[]
  created_at: string
}

export interface ActivityPayload {
  title: string
  description?: string
  difficulty: string
  duration_minutes: number
  equipment_requirement?: string
  start_time: string
  end_time: string
  start_lat: number
  start_lng: number
  end_lat: number
  end_lng: number
  address: string
  max_teams: number
}

export function listActivities(params: { page?: number; page_size?: number; status?: string; difficulty?: string; keyword?: string }) {
  return request.get<{ list: Activity[]; total: number; page: number; page_size: number }>('/activities', { params })
}

export function getActivity(id: number) {
  return request.get<Activity>(`/activities/${id}`)
}

export function createActivity(data: ActivityPayload) {
  return request.post('/activities', data)
}

export function updateActivity(id: number, data: ActivityPayload) {
  return request.put(`/activities/${id}`, data)
}

export function transitionActivity(id: number, status: string) {
  return request.post(`/activities/${id}/transition`, { status })
}

export function deleteActivity(id: number) {
  return request.delete(`/activities/${id}`)
}
