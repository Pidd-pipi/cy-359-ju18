import { create } from 'zustand'
import * as checkinApi from '../api/checkin'
import type { LeaderboardRow } from '../api/checkin'

interface CheckinState {
  leaderboard: Record<number, LeaderboardRow[]>
  fetchLeaderboard: (activityId: number) => Promise<void>
}

export const useCheckinStore = create<CheckinState>((set) => ({
  leaderboard: {},
  fetchLeaderboard: async (activityId) => {
    const res = await checkinApi.getLeaderboard(activityId)
    set((s) => ({ leaderboard: { ...s.leaderboard, [activityId]: res.data } }))
  },
}))
