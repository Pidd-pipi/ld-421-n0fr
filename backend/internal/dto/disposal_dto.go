package dto

import "time"

// SubmitDisposalRequest 提交设备报废处置申请请求。
type SubmitDisposalRequest struct {
	Reason string `json:"reason" binding:"required,max=512"`
}

// DisposalBlockerResponse 处置审批阻塞项。
type DisposalBlockerResponse struct {
	Type         string `json:"type"`
	RecordID     uint   `json:"recordId"`
	UserID       uint   `json:"userId"`
	UserName     string `json:"userName"`
	Status       string `json:"status"`
	Detail       string `json:"detail"`
	ExpectedTime string `json:"expectedTime,omitempty"`
}

// DisposalResponse 处置审批信息返回。
type DisposalResponse struct {
	ID          uint                      `json:"id"`
	EquipmentID uint                      `json:"equipmentId"`
	Reason      string                    `json:"reason"`
	Status      string                    `json:"status"`
	ApplicantID uint                      `json:"applicantId"`
	Applicant   string                    `json:"applicantName,omitempty"`
	ReviewerID  *uint                     `json:"reviewerId"`
	Reviewer    string                    `json:"reviewerName,omitempty"`
	Blockers    []DisposalBlockerResponse `json:"blockers"`
	ReviewedAt  *time.Time                `json:"reviewedAt"`
	Active      bool                      `json:"active"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
}

// DisposalReviewResponse 处置审批结果返回，退回时携带阻塞项。
type DisposalReviewResponse struct {
	Approval DisposalResponse `json:"approval"`
	Approved bool             `json:"approved"`
	Blockers []DisposalBlockerResponse `json:"blockers"`
}
