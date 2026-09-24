package model

import (
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
)

// BorrowRecord 设备借用记录实体。
type BorrowRecord struct {
	Base
	EquipmentID        uint                       `gorm:"index;not null" json:"equipmentId"`
	BorrowerID         uint                       `gorm:"index;not null" json:"borrowerId"`
	BorrowDate         time.Time                  `json:"borrowDate"`
	ExpectedReturnDate time.Time                  `json:"expectedReturnDate"`
	ActualReturnDate   *time.Time                 `json:"actualReturnDate"`
	Reason             string                     `gorm:"size:512" json:"reason"`
	Status             constants.BorrowStatus     `gorm:"size:32;index;not null;default:Pending" json:"status"`
	ApproverID         *uint                      `json:"approverId"`
	ReturnCondition    *constants.ReturnCondition `gorm:"size:32" json:"returnCondition"`
	Equipment          *Equipment                 `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Borrower           *User                      `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
	Approver           *User                      `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

func (BorrowRecord) TableName() string { return "borrow_records" }
