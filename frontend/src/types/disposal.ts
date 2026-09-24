export interface DisposalBlocker {
  type: 'borrow' | 'reservation' | string
  recordId: number
  status: string
  userName: string
  detail: string
}

export interface DisposalRequest {
  id: number
  equipmentId: number
  equipmentName?: string
  reason: string
  status: string
  applicantId: number
  applicantName?: string
  approverId?: number
  approverName?: string
  blockers: DisposalBlocker[]
  processedAt?: string
  createdAt: string
}
