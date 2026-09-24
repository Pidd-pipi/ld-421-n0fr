package dto

import "time"

// SubmitDisposalRequest 提交报废申请请求。
type SubmitDisposalRequest struct {
	Reason string `json:"reason" binding:"required,max=512"`
}

// DisposalBlocker 报废审批阻塞项返回。
type DisposalBlocker struct {
	Type     string `json:"type"`
	RecordID uint   `json:"recordId"`
	Status   string `json:"status"`
	UserName string `json:"userName"`
	Detail   string `json:"detail"`
}

// DisposalResponse 报废申请返回。
type DisposalResponse struct {
	ID            uint              `json:"id"`
	EquipmentID   uint              `json:"equipmentId"`
	EquipmentName string            `json:"equipmentName,omitempty"`
	Reason        string            `json:"reason"`
	Status        string            `json:"status"`
	ApplicantID   uint              `json:"applicantId"`
	ApplicantName string            `json:"applicantName,omitempty"`
	ApproverID    *uint             `json:"approverId"`
	ApproverName  string            `json:"approverName,omitempty"`
	Blockers      []DisposalBlocker `json:"blockers"`
	ProcessedAt   *time.Time        `json:"processedAt"`
	CreatedAt     time.Time         `json:"createdAt"`
}
