package constants

// AssetStatus 设备资产状态。
type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "Available"
	AssetStatusInUse       AssetStatus = "InUse"
	AssetStatusMaintenance AssetStatus = "Maintenance"
	AssetStatusRetired     AssetStatus = "Retired"
	AssetStatusLost        AssetStatus = "Lost"
)

// AllAssetStatus 返回全部资产状态，供校验与下拉使用。
func AllAssetStatus() []AssetStatus {
	return []AssetStatus{
		AssetStatusAvailable,
		AssetStatusInUse,
		AssetStatusMaintenance,
		AssetStatusRetired,
		AssetStatusLost,
	}
}

func (s AssetStatus) Valid() bool {
	for _, v := range AllAssetStatus() {
		if v == s {
			return true
		}
	}
	return false
}

// BorrowStatus 借用审批状态。
type BorrowStatus string

const (
	BorrowStatusPending  BorrowStatus = "Pending"
	BorrowStatusApproved BorrowStatus = "Approved"
	BorrowStatusRejected BorrowStatus = "Rejected"
	BorrowStatusReturned BorrowStatus = "Returned"
	BorrowStatusOverdue  BorrowStatus = "Overdue"
)

func AllBorrowStatus() []BorrowStatus {
	return []BorrowStatus{
		BorrowStatusPending,
		BorrowStatusApproved,
		BorrowStatusRejected,
		BorrowStatusReturned,
		BorrowStatusOverdue,
	}
}

func (s BorrowStatus) Valid() bool {
	for _, v := range AllBorrowStatus() {
		if v == s {
			return true
		}
	}
	return false
}

// MaintenanceType 维护类型。
type MaintenanceType string

const (
	MaintenanceTypePreventive  MaintenanceType = "Preventive"
	MaintenanceTypeCorrective  MaintenanceType = "Corrective"
	MaintenanceTypeCalibration MaintenanceType = "Calibration"
	MaintenanceTypeCleaning    MaintenanceType = "Cleaning"
)

func AllMaintenanceType() []MaintenanceType {
	return []MaintenanceType{
		MaintenanceTypePreventive,
		MaintenanceTypeCorrective,
		MaintenanceTypeCalibration,
		MaintenanceTypeCleaning,
	}
}

func (t MaintenanceType) Valid() bool {
	for _, v := range AllMaintenanceType() {
		if v == t {
			return true
		}
	}
	return false
}

// ReturnCondition 归还状况。
type ReturnCondition string

const (
	ReturnConditionGood    ReturnCondition = "Good"
	ReturnConditionDamaged ReturnCondition = "Damaged"
	ReturnConditionLost    ReturnCondition = "Lost"
)

func AllReturnCondition() []ReturnCondition {
	return []ReturnCondition{ReturnConditionGood, ReturnConditionDamaged, ReturnConditionLost}
}

func (c ReturnCondition) Valid() bool {
	for _, v := range AllReturnCondition() {
		if v == c {
			return true
		}
	}
	return false
}

// MaintenanceResult 维护结果。
type MaintenanceResult string

const (
	MaintenanceResultPass          MaintenanceResult = "Pass"
	MaintenanceResultFail          MaintenanceResult = "Fail"
	MaintenanceResultNeedsFollowUp MaintenanceResult = "NeedsFollowUp"
)

func AllMaintenanceResult() []MaintenanceResult {
	return []MaintenanceResult{
		MaintenanceResultPass,
		MaintenanceResultFail,
		MaintenanceResultNeedsFollowUp,
	}
}

func (r MaintenanceResult) Valid() bool {
	for _, v := range AllMaintenanceResult() {
		if v == r {
			return true
		}
	}
	return false
}

// ReservationStatus 预约审批状态。
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "Pending"
	ReservationStatusApproved  ReservationStatus = "Approved"
	ReservationStatusRejected  ReservationStatus = "Rejected"
	ReservationStatusCancelled ReservationStatus = "Cancelled"
)

func AllReservationStatus() []ReservationStatus {
	return []ReservationStatus{
		ReservationStatusPending,
		ReservationStatusApproved,
		ReservationStatusRejected,
		ReservationStatusCancelled,
	}
}

func (s ReservationStatus) Valid() bool {
	for _, v := range AllReservationStatus() {
		if v == s {
			return true
		}
	}
	return false
}
