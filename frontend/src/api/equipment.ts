import { get, post, put, patch } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { Equipment, EquipmentPayload, EquipmentQuery, PageResult } from '../types'

export function fetchEquipment(params?: EquipmentQuery) {
  return get<PageResult<Equipment>>(API_PATHS.equipment, params as Record<string, unknown>)
}

export function fetchEquipmentDetail(id: number | string) {
  return get<Equipment>(API_PATHS.equipmentDetail(id))
}

export function createEquipment(payload: EquipmentPayload) {
  return post<Equipment>(API_PATHS.equipment, payload)
}

export function updateEquipment(id: number | string, payload: EquipmentPayload) {
  return put<Equipment>(API_PATHS.equipmentDetail(id), payload)
}

export function retireEquipment(id: number | string) {
  return patch<null>(API_PATHS.equipmentStatus(id), { status: 'Retired' })
}

export function transferOwner(id: number | string, ownerId: number) {
  return patch<null>(API_PATHS.equipmentOwner(id), { ownerId })
}
