export interface Equipment {
  id: number
  name: string
  code: string
  categoryId: number
  categoryName?: string
  brandModel: string
  serialNumber: string
  purchaseDate?: string
  purchasePrice: number
  location: string
  status: string
  ownerId: number
  ownerName?: string
  supplier: string
  warrantyExpiry?: string
  imageUrl: string
  createdAt: string
  updatedAt: string
}

export interface EquipmentPayload {
  name: string
  code: string
  categoryId: number
  brandModel?: string
  serialNumber?: string
  purchaseDate?: string
  purchasePrice?: number
  location?: string
  ownerId: number
  supplier?: string
  warrantyExpiry?: string
  imageUrl?: string
}

export interface EquipmentQuery {
  page?: number
  page_size?: number
  keyword?: string
  category_id?: number
  status?: string
}
