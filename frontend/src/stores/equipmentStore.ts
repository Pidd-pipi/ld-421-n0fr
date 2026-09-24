import { create } from 'zustand'
import { fetchEquipment } from '../api/equipment'
import type { Equipment, EquipmentQuery } from '../types'

interface EquipmentState {
  items: Equipment[]
  total: number
  loading: boolean
  fetch: (params?: EquipmentQuery) => Promise<void>
}

export const useEquipmentStore = create<EquipmentState>((set) => ({
  items: [],
  total: 0,
  loading: false,
  fetch: async (params) => {
    set({ loading: true })
    try {
      const result = await fetchEquipment(params)
      set({ items: result.list, total: result.total, loading: false })
    } catch (error) {
      set({ loading: false })
      throw error
    }
  }
}))
