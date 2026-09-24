import { get, post } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { DisposalApproval, DisposalReviewResult, PageResult, SubmitDisposalPayload } from '../types'

export function fetchDisposals(params?: Record<string, unknown>) {
  return get<PageResult<DisposalApproval>>(API_PATHS.disposals, params)
}

export function fetchDisposalDetail(id: number | string) {
  return get<DisposalApproval>(API_PATHS.disposalDetail(id))
}

export function fetchEquipmentDisposal(equipmentId: number | string) {
  return get<DisposalApproval | null>(API_PATHS.equipmentDisposal(equipmentId))
}

export function submitDisposal(equipmentId: number | string, payload: SubmitDisposalPayload) {
  return post<DisposalApproval>(API_PATHS.equipmentDisposal(equipmentId), payload)
}

export function reviewDisposal(id: number | string) {
  return post<DisposalReviewResult>(API_PATHS.disposalReview(id))
}
