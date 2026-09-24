export const AssetStatus = {
  Available: 'Available',
  InUse: 'InUse',
  Maintenance: 'Maintenance',
  Retired: 'Retired',
  Lost: 'Lost'
} as const
export type AssetStatus = (typeof AssetStatus)[keyof typeof AssetStatus]

export const BorrowStatus = {
  Pending: 'Pending',
  Approved: 'Approved',
  Rejected: 'Rejected',
  Returned: 'Returned',
  Overdue: 'Overdue'
} as const
export type BorrowStatus = (typeof BorrowStatus)[keyof typeof BorrowStatus]

export const MaintenanceType = {
  Preventive: 'Preventive',
  Corrective: 'Corrective',
  Calibration: 'Calibration',
  Cleaning: 'Cleaning'
} as const
export type MaintenanceType = (typeof MaintenanceType)[keyof typeof MaintenanceType]

export const ReturnCondition = {
  Good: 'Good',
  Damaged: 'Damaged',
  Lost: 'Lost'
} as const
export type ReturnCondition = (typeof ReturnCondition)[keyof typeof ReturnCondition]

export const MaintenanceResult = {
  Pass: 'Pass',
  Fail: 'Fail',
  NeedsFollowUp: 'NeedsFollowUp'
} as const
export type MaintenanceResult = (typeof MaintenanceResult)[keyof typeof MaintenanceResult]

export const ReservationStatus = {
  Pending: 'Pending',
  Approved: 'Approved',
  Rejected: 'Rejected',
  Cancelled: 'Cancelled'
} as const
export type ReservationStatus = (typeof ReservationStatus)[keyof typeof ReservationStatus]
