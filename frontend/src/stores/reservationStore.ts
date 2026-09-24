import { create } from 'zustand'
import { fetchReservations } from '../api/reservation'
import type { Reservation } from '../types'

interface ReservationState {
  items: Reservation[]
  total: number
  loading: boolean
  fetch: (params?: Record<string, unknown>) => Promise<void>
}

export const useReservationStore = create<ReservationState>((set) => ({
  items: [],
  total: 0,
  loading: false,
  fetch: async (params) => {
    set({ loading: true })
    try {
      const result = await fetchReservations(params)
      set({ items: result.list, total: result.total, loading: false })
    } catch (error) {
      set({ loading: false })
      throw error
    }
  }
}))
