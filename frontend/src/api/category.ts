import { get, post, put, del } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { CategoryPayload, EquipmentCategory } from '../types'

export function fetchCategories() {
  return get<EquipmentCategory[]>(API_PATHS.categories)
}

export function createCategory(payload: CategoryPayload) {
  return post<EquipmentCategory>(API_PATHS.categories, payload)
}

export function updateCategory(id: number | string, payload: CategoryPayload) {
  return put<EquipmentCategory>(`${API_PATHS.categories}/${id}`, payload)
}

export function deleteCategory(id: number | string) {
  return del<null>(`${API_PATHS.categories}/${id}`)
}
