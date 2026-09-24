export interface EquipmentCategory {
  id: number
  name: string
  parentId?: number
  description: string
  icon: string
  children?: EquipmentCategory[]
}

export interface CategoryPayload {
  name: string
  parentId?: number
  description?: string
  icon?: string
}
