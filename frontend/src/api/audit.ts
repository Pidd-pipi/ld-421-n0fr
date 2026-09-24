import { get } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'

export interface AuditLog {
  id: number
  userId: number
  userName: string
  action: string
  resourceType: string
  resourceId: number
  detail: string
  ip: string
  createdAt: string
}

export function fetchAuditLogs(params?: Record<string, unknown>) {
  return get<{ list: AuditLog[]; total: number; page: number; pageSize: number }>(API_PATHS.auditLogs, params)
}
