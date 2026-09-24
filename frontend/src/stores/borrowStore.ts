import { create } from 'zustand'
import { fetchBorrows } from '../api/borrow'
import type { BorrowRecord } from '../types'

interface BorrowState {
  items: BorrowRecord[]
  total: number
  loading: boolean
  fetch: (params?: Record<string, unknown>) => Promise<void>
}

export const useBorrowStore = create<BorrowState>((set) => ({
  items: [],
  total: 0,
  loading: false,
  fetch: async (params) => {
    set({ loading: true })
    try {
      const result = await fetchBorrows(params)
      set({ items: result.list, total: result.total, loading: false })
    } catch (error) {
      set({ loading: false })
      throw error
    }
  }
}))
