import { get, post } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { CreateReservationPayload, PageResult, Reservation } from '../types'

export function fetchReservations(params?: Record<string, unknown>) {
  return get<PageResult<Reservation>>(API_PATHS.reservations, params)
}

export function fetchReservationDetail(id: number | string) {
  return get<Reservation>(API_PATHS.reservationDetail(id))
}

export function createReservation(payload: CreateReservationPayload) {
  return post<Reservation>(API_PATHS.reservations, payload)
}

export function approveReservation(id: number | string) {
  return post<null>(API_PATHS.reservationApprove(id))
}

export function rejectReservation(id: number | string) {
  return post<null>(API_PATHS.reservationReject(id))
}

export function cancelReservation(id: number | string) {
  return post<null>(API_PATHS.reservationCancel(id))
}
