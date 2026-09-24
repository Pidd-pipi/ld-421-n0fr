import { get, post } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { DisposalRequest } from '../types'

export function fetchDisposals(equipmentId: number | string) {
  return get<DisposalRequest[]>(API_PATHS.equipmentDisposals(equipmentId))
}

export function submitDisposal(equipmentId: number | string, reason: string) {
  return post<DisposalRequest>(API_PATHS.equipmentDisposals(equipmentId), { reason })
}

export function approveDisposal(equipmentId: number | string, disposalId: number | string) {
  return post<DisposalRequest>(API_PATHS.disposalApprove(equipmentId, disposalId))
}
