export interface BorrowRecord {
  id: number
  equipmentId: number
  equipmentName?: string
  equipmentCode?: string
  borrowerId: number
  borrowerName?: string
  borrowDate: string
  expectedReturnDate: string
  actualReturnDate?: string
  reason: string
  status: string
  approverId?: number
  approverName?: string
  returnCondition?: string
  createdAt: string
}

export interface CreateBorrowPayload {
  equipmentId: number
  borrowDate: string
  expectedReturnDate: string
  reason?: string
}

export interface ReturnBorrowPayload {
  actualReturnDate: string
  returnCondition: string
}
