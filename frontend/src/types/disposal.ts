import type { DisposalStatus } from './enums'

export interface DisposalBlocker {
  type: 'Borrow' | 'Reservation' | string
  recordId: number
  userId: number
  userName: string
  status: string
  detail: string
  expectedTime?: string
}

export interface DisposalApproval {
  id: number
  equipmentId: number
  reason: string
  status: DisposalStatus | string
  applicantId: number
  applicantName?: string
  reviewerId?: number
  reviewerName?: string
  blockers: DisposalBlocker[]
  reviewedAt?: string
  active: boolean
  createdAt: string
  updatedAt: string
}

export interface DisposalReviewResult {
  approval: DisposalApproval
  approved: boolean
  blockers: DisposalBlocker[]
}

export interface SubmitDisposalPayload {
  reason: string
}
