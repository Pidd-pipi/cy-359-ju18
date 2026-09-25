import request from '../utils/request'

export interface CheckinRecord {
  id: number
  activity_id: number
  checkpoint_id: number
  team_id: number
  user_id: number
  checkin_type: string
  latitude: number
  longitude: number
  answer?: string
  photo_url?: string
  result: string
  points_earned: number
  checked_in_at: string
}

export interface LeaderboardRow {
  rank: number
  team_id: number
  team_name: string
  total_seconds: number
  duration: string
  checkpoint_count: number
  status: string
}

export function checkin(teamId: number, data: { checkpoint_id: number; checkin_type: string; latitude?: number; longitude?: number; answer?: string; photo_url?: string }) {
  return request.post(`/teams/${teamId}/checkin`, data)
}

export function listTeamCheckins(teamId: number, activityId?: number) {
  return request.get<CheckinRecord[]>(`/teams/${teamId}/checkins`, { params: { activity_id: activityId } })
}

export function listActivityCheckins(activityId: number) {
  return request.get<CheckinRecord[]>(`/activities/${activityId}/checkins`)
}

export function getLeaderboard(activityId: number) {
  return request.get<LeaderboardRow[]>(`/activities/${activityId}/leaderboard`)
}
