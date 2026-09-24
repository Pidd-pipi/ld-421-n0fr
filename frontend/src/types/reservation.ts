export interface Reservation {
  id: number
  equipmentId: number
  equipmentName?: string
  userId: number
  userName?: string
  startTime: string
  endTime: string
  purpose: string
  status: string
  approverId?: number
  approverName?: string
  createdAt: string
}

export interface CreateReservationPayload {
  equipmentId: number
  startTime: string
  endTime: string
  purpose?: string
}
