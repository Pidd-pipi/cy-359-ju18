import { create } from 'zustand'
import * as auditApi from '../api/audit'
import type { AuditLog } from '../api/audit'

interface AuditState {
  logs: AuditLog[]
  total: number
  fetchLogs: (params?: { page?: number; page_size?: number; action?: string; resource_type?: string; keyword?: string }) => Promise<void>
}

export const useAuditStore = create<AuditState>((set) => ({
  logs: [],
  total: 0,
  fetchLogs: async (params = {}) => {
    const res = await auditApi.listAuditLogs({ page: 1, page_size: 20, ...params })
    set({ logs: res.data.list, total: res.data.total })
  },
}))
