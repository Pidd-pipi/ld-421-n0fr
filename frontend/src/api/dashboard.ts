import { get } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'

export interface DashboardStats {
  statusDistribution: Record<string, number>
  topBorrows: { equipmentId: number; name: string; code: string; count: number }[]
  expiringWarranty: { id: number; name: string; code: string; warrantyExpiry?: string }[]
  pendingBorrows: number
  pendingReservations: number
}

export function fetchDashboardStats() {
  return get<DashboardStats>(API_PATHS.dashboardStats)
}
