import request from '../utils/request'

export interface AuditLog {
  id: number
  user_id: number
  username: string
  action: string
  resource_type: string
  resource_id: string
  detail: string
  ip: string
  request_id: string
  created_at: string
}

export function listAuditLogs(params: { page?: number; page_size?: number; action?: string; resource_type?: string; keyword?: string }) {
  return request.get<{ list: AuditLog[]; total: number }>('/audit-logs', { params })
}
