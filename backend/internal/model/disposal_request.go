package model

import (
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
)

// DisposalBlocker 报废审批阻塞项，描述仍未完结的借用或预约。
type DisposalBlocker struct {
	Type     string `json:"type"` // borrow / reservation
	RecordID uint   `json:"recordId"`
	Status   string `json:"status"`
	UserName string `json:"userName"`
	Detail   string `json:"detail"`
}

// DisposalRequest 设备报废处置审批单。审批通过后设备才转为 Retired。
type DisposalRequest struct {
	Base
	EquipmentID uint                     `gorm:"index;not null" json:"equipmentId"`
	Reason      string                   `gorm:"size:512;not null" json:"reason"`
	Status      constants.DisposalStatus `gorm:"size:32;index;not null;default:Pending" json:"status"`
	ApplicantID uint                     `gorm:"index;not null" json:"applicantId"`
	ApproverID  *uint                    `json:"approverId"`
	Blockers    string                   `gorm:"type:text" json:"-"` // 退回时的阻塞项 JSON 快照
	ProcessedAt *time.Time               `json:"processedAt"`
	Equipment   *Equipment               `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Applicant   *User                    `gorm:"foreignKey:ApplicantID" json:"applicant,omitempty"`
	Approver    *User                    `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

func (DisposalRequest) TableName() string { return "disposal_requests" }
