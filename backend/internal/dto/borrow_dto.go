package dto

import "time"

// CreateBorrowRequest 提交借用申请。
type CreateBorrowRequest struct {
	EquipmentID        uint   `json:"equipmentId" binding:"required"`
	BorrowDate         string `json:"borrowDate" binding:"required"`
	ExpectedReturnDate string `json:"expectedReturnDate" binding:"required"`
	Reason             string `json:"reason" binding:"omitempty,max=512"`
}

// ReturnBorrowRequest 确认归还。
type ReturnBorrowRequest struct {
	ActualReturnDate string `json:"actualReturnDate" binding:"required"`
	ReturnCondition  string `json:"returnCondition" binding:"required"`
}

// BorrowResponse 借用记录返回。
type BorrowResponse struct {
	ID                 uint       `json:"id"`
	EquipmentID        uint       `json:"equipmentId"`
	EquipmentName      string     `json:"equipmentName,omitempty"`
	EquipmentCode      string     `json:"equipmentCode,omitempty"`
	BorrowerID         uint       `json:"borrowerId"`
	BorrowerName       string     `json:"borrowerName,omitempty"`
	BorrowDate         time.Time  `json:"borrowDate"`
	ExpectedReturnDate time.Time  `json:"expectedReturnDate"`
	ActualReturnDate   *time.Time `json:"actualReturnDate"`
	Reason             string     `json:"reason"`
	Status             string     `json:"status"`
	ApproverID         *uint      `json:"approverId"`
	ApproverName       string     `json:"approverName,omitempty"`
	ReturnCondition    *string    `json:"returnCondition"`
	CreatedAt          time.Time  `json:"createdAt"`
}
