import { useState } from 'react'
import { approveBorrow, createBorrow, rejectBorrow, returnBorrow } from '../api/borrow'
import type { CreateBorrowPayload, ReturnBorrowPayload } from '../types'

const STATUS_STEP: Record<string, number> = {
  Pending: 0,
  Approved: 1,
  Returned: 3,
  Rejected: 1,
  Overdue: 2
}

export function useBorrowFlow() {
  const [loading, setLoading] = useState(false)

  const submit = async (payload: CreateBorrowPayload) => {
    setLoading(true)
    try {
      const result = await createBorrow(payload)
      return result
    } finally {
      setLoading(false)
    }
  }

  const approve = async (id: number | string) => {
    setLoading(true)
    try {
      await approveBorrow(id)
    } finally {
      setLoading(false)
    }
  }

  const reject = async (id: number | string) => {
    setLoading(true)
    try {
      await rejectBorrow(id)
    } finally {
      setLoading(false)
    }
  }

  const confirmReturn = async (id: number | string, payload: ReturnBorrowPayload) => {
    setLoading(true)
    try {
      await returnBorrow(id, payload)
    } finally {
      setLoading(false)
    }
  }

  const stepOf = (status: string) => STATUS_STEP[status] ?? 0

  return { loading, submit, approve, reject, confirmReturn, stepOf }
}

export default useBorrowFlow
