import request from '../utils/request'
import type { Checkpoint } from './activity'

export interface CheckpointPayload {
  name: string
  sequence: number
  lat: number
  lng: number
  clue?: string
  task_type: string
  task_content?: string
  expected_answer?: string
  radius_meters: number
  qr_code: string
}

export function createCheckpoint(activityId: number, data: CheckpointPayload) {
  return request.post(`/checkpoints/activities/${activityId}`, data)
}

export function updateCheckpoint(id: number, data: CheckpointPayload) {
  return request.put(`/checkpoints/${id}`, data)
}

export function deleteCheckpoint(id: number) {
  return request.delete(`/checkpoints/${id}`)
}

export function listCheckpoints(activityId: number) {
  return request.get<Checkpoint[]>(`/activities/${activityId}`)
}
