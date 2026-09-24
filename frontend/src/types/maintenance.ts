export interface MaintenanceRecord {
  id: number
  equipmentId: number
  equipmentName?: string
  type: string
  content: string
  maintenanceDate: string
  nextMaintenanceDate?: string
  cost: number
  maintainerId: number
  maintainerName?: string
  result: string
  createdAt: string
}

export interface CreateMaintenancePayload {
  equipmentId: number
  type: string
  content?: string
  maintenanceDate: string
  nextMaintenanceDate?: string
  cost?: number
  maintainerId: number
}

export interface MaintenanceStats {
  totalRecords: number
  totalCost: number
  resultCount: Record<string, number>
}
