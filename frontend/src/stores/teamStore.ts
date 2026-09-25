import { create } from 'zustand'
import * as teamApi from '../api/team'
import type { Team, Registration } from '../api/team'

interface TeamState {
  myTeams: Team[]
  myRegistrations: Registration[]
  fetchMyTeams: () => Promise<void>
  fetchMyRegistrations: () => Promise<void>
}

export const useTeamStore = create<TeamState>((set) => ({
  myTeams: [],
  myRegistrations: [],
  fetchMyTeams: async () => {
    const res = await teamApi.listMyTeams()
    set({ myTeams: res.data })
  },
  fetchMyRegistrations: async () => {
    const res = await teamApi.listMyRegistrations()
    set({ myRegistrations: res.data })
  },
}))
