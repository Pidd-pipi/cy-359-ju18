import { create } from 'zustand'
import * as activityApi from '../api/activity'
import type { Activity, ActivityPayload } from '../api/activity'

interface ActivityState {
  list: Activity[]
  total: number
  loading: boolean
  fetchList: (params?: { page?: number; page_size?: number; status?: string; difficulty?: string; keyword?: string }) => Promise<void>
}

export const useActivityStore = create<ActivityState>((set) => ({
  list: [],
  total: 0,
  loading: false,
  fetchList: async (params = {}) => {
    set({ loading: true })
    try {
      const res = await activityApi.listActivities({ page: 1, page_size: 12, ...params })
      set({ list: res.data.list, total: res.data.total })
    } finally {
      set({ loading: false })
    }
  },
}))

export type { Activity, ActivityPayload }
