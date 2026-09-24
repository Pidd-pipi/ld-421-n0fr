import { create } from 'zustand'
import { fetchMaintenance } from '../api/maintenance'
import type { MaintenanceRecord } from '../types'

interface MaintenanceState {
  items: MaintenanceRecord[]
  total: number
  loading: boolean
  fetch: (params?: Record<string, unknown>) => Promise<void>
}

export const useMaintenanceStore = create<MaintenanceState>((set) => ({
  items: [],
  total: 0,
  loading: false,
  fetch: async (params) => {
    set({ loading: true })
    try {
      const result = await fetchMaintenance(params)
      set({ items: result.list, total: result.total, loading: false })
    } catch (error) {
      set({ loading: false })
      throw error
    }
  }
}))
