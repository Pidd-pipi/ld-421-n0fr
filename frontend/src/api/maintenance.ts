import { get, post } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { CreateMaintenancePayload, MaintenanceRecord, MaintenanceStats, PageResult } from '../types'

export function fetchMaintenance(params?: Record<string, unknown>) {
  return get<PageResult<MaintenanceRecord>>(API_PATHS.maintenance, params)
}

export function createMaintenance(payload: CreateMaintenancePayload) {
  return post<MaintenanceRecord>(API_PATHS.maintenance, payload)
}

export function executeMaintenance(id: number | string, result: string) {
  return post<null>(API_PATHS.maintenanceExecute(id), { result })
}

export function fetchMaintenanceStats() {
  return get<MaintenanceStats>(API_PATHS.maintenanceStats)
}
