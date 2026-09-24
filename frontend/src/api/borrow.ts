import { get, post } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { BorrowRecord, CreateBorrowPayload, PageResult, ReturnBorrowPayload } from '../types'

export function fetchBorrows(params?: Record<string, unknown>) {
  return get<PageResult<BorrowRecord>>(API_PATHS.borrows, params)
}

export function fetchBorrowDetail(id: number | string) {
  return get<BorrowRecord>(API_PATHS.borrowDetail(id))
}

export function createBorrow(payload: CreateBorrowPayload) {
  return post<BorrowRecord>(API_PATHS.borrows, payload)
}

export function approveBorrow(id: number | string) {
  return post<null>(API_PATHS.borrowApprove(id))
}

export function rejectBorrow(id: number | string) {
  return post<null>(API_PATHS.borrowReject(id))
}

export function returnBorrow(id: number | string, payload: ReturnBorrowPayload) {
  return post<null>(API_PATHS.borrowReturn(id), payload)
}
