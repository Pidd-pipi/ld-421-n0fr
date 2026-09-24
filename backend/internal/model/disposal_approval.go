package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
)

// DisposalBlocker 处置审批阻塞项（未归还借用 / 未结束预约）。
type DisposalBlocker struct {
	Type         string `json:"type"`
	RecordID     uint   `json:"recordId"`
	UserID       uint   `json:"userId"`
	UserName     string `json:"userName"`
	Status       string `json:"status"`
	Detail       string `json:"detail"`
	ExpectedTime string `json:"expectedTime,omitempty"`
}

// DisposalBlockers 阻塞项列表，以 JSON 形式持久化。
type DisposalBlockers []DisposalBlocker

// Value 实现 driver.Valuer。
func (b DisposalBlockers) Value() (driver.Value, error) {
	if b == nil {
		return "[]", nil
	}
	data, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("marshal disposal blockers: %w", err)
	}
	return string(data), nil
}

// Scan 实现 sql.Scanner。
func (b *DisposalBlockers) Scan(src any) error {
	if src == nil {
		*b = DisposalBlockers{}
		return nil
	}
	var data []byte
	switch value := src.(type) {
	case []byte:
		data = value
	case string:
		data = []byte(value)
	default:
		return errors.New("unsupported disposal blockers source type")
	}
	if len(data) == 0 {
		*b = DisposalBlockers{}
		return nil
	}
	var result DisposalBlockers
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("unmarshal disposal blockers: %w", err)
	}
	*b = result
	return nil
}

// DisposalApproval 设备处置（报废）审批单。
//
// 审批进行中时 ActiveEquipmentID 等于 EquipmentID，审批结束后置为 NULL；
// 借助唯一索引中 NULL 不冲突的特性，保证同一设备最多只有一个进行中的审批，
// 同时保留全部历史审批记录。
type DisposalApproval struct {
	Base
	EquipmentID       uint                     `gorm:"index;not null" json:"equipmentId"`
	Reason            string                   `gorm:"size:512;not null" json:"reason"`
	Status            constants.DisposalStatus `gorm:"size:32;index;not null;default:Pending" json:"status"`
	ApplicantID       uint                     `gorm:"index;not null" json:"applicantId"`
	ReviewerID        *uint                    `json:"reviewerId"`
	Blockers          DisposalBlockers         `gorm:"type:json" json:"blockers"`
	ReviewedAt        *time.Time               `json:"reviewedAt"`
	Active            bool                     `gorm:"index;not null;default:true" json:"active"`
	ActiveEquipmentID *uint                    `gorm:"column:active_equipment_id;uniqueIndex" json:"-"`
	Equipment         *Equipment               `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Applicant         *User                    `gorm:"foreignKey:ApplicantID" json:"applicant,omitempty"`
	Reviewer          *User                    `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// TableName 指定处置审批表名。
func (DisposalApproval) TableName() string { return "disposal_approvals" }
