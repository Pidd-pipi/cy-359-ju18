import request from '../utils/request'
import type { Activity } from './activity'

export interface Favorite {
  id: number
  activity_id: number
  activity: Activity
  created_at: string
}

export function listFavorites(params: { page?: number; page_size?: number }) {
  return request.get<{ list: Favorite[]; total: number }>('/favorites', { params })
}

export function addFavorite(activityId: number) {
  return request.post('/favorites', { activity_id: activityId })
}

export function removeFavorite(activityId: number) {
  return request.delete(`/favorites/${activityId}`)
}
