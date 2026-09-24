package handler

import (
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/service"
)

func mapUser(user model.User) dto.UserResponse {
	roleCode := ""
	roleName := ""
	if user.Role != nil {
		roleCode = user.Role.Code
		roleName = user.Role.Name
	}
	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		RoleID:   user.RoleID,
		RoleCode: roleCode,
		RoleName: roleName,
	}
}

func mapEquipment(equipment model.Equipment) dto.EquipmentResponse {
	categoryName := ""
	if equipment.Category != nil {
		categoryName = equipment.Category.Name
	}
	ownerName := ""
	if equipment.Owner != nil {
		ownerName = equipment.Owner.Name
	}
	return dto.EquipmentResponse{
		ID:             equipment.ID,
		Name:           equipment.Name,
		Code:           equipment.Code,
		CategoryID:     equipment.CategoryID,
		CategoryName:   categoryName,
		BrandModel:     equipment.BrandModel,
		SerialNumber:   equipment.SerialNumber,
		PurchaseDate:   equipment.PurchaseDate,
		PurchasePrice:  equipment.PurchasePrice,
		Location:       equipment.Location,
		Status:         string(equipment.Status),
		OwnerID:        equipment.OwnerID,
		OwnerName:      ownerName,
		Supplier:       equipment.Supplier,
		WarrantyExpiry: equipment.WarrantyExpiry,
		ImageURL:       equipment.ImageURL,
		CreatedAt:      equipment.CreatedAt,
		UpdatedAt:      equipment.UpdatedAt,
	}
}

func mapBorrow(record model.BorrowRecord) dto.BorrowResponse {
	equipmentName := ""
	equipmentCode := ""
	if record.Equipment != nil {
		equipmentName = record.Equipment.Name
		equipmentCode = record.Equipment.Code
	}
	borrowerName := ""
	if record.Borrower != nil {
		borrowerName = record.Borrower.Name
	}
	approverName := ""
	if record.Approver != nil {
		approverName = record.Approver.Name
	}
	var returnCondition *string
	if record.ReturnCondition != nil {
		value := string(*record.ReturnCondition)
		returnCondition = &value
	}
	return dto.BorrowResponse{
		ID:                 record.ID,
		EquipmentID:        record.EquipmentID,
		EquipmentName:      equipmentName,
		EquipmentCode:      equipmentCode,
		BorrowerID:         record.BorrowerID,
		BorrowerName:       borrowerName,
		BorrowDate:         record.BorrowDate,
		ExpectedReturnDate: record.ExpectedReturnDate,
		ActualReturnDate:   record.ActualReturnDate,
		Reason:             record.Reason,
		Status:             string(record.Status),
		ApproverID:         record.ApproverID,
		ApproverName:       approverName,
		ReturnCondition:    returnCondition,
		CreatedAt:          record.CreatedAt,
	}
}

func mapMaintenance(record model.MaintenanceRecord) dto.MaintenanceResponse {
	equipmentName := ""
	if record.Equipment != nil {
		equipmentName = record.Equipment.Name
	}
	maintainerName := ""
	if record.Maintainer != nil {
		maintainerName = record.Maintainer.Name
	}
	return dto.MaintenanceResponse{
		ID:                  record.ID,
		EquipmentID:         record.EquipmentID,
		EquipmentName:       equipmentName,
		Type:                string(record.Type),
		Content:             record.Content,
		MaintenanceDate:     record.MaintenanceDate,
		NextMaintenanceDate: record.NextMaintenanceDate,
		Cost:                record.Cost,
		MaintainerID:        record.MaintainerID,
		MaintainerName:      maintainerName,
		Result:              string(record.Result),
		CreatedAt:           record.CreatedAt,
	}
}

func mapReservation(reservation model.Reservation) dto.ReservationResponse {
	equipmentName := ""
	if reservation.Equipment != nil {
		equipmentName = reservation.Equipment.Name
	}
	userName := ""
	if reservation.User != nil {
		userName = reservation.User.Name
	}
	approverName := ""
	if reservation.Approver != nil {
		approverName = reservation.Approver.Name
	}
	return dto.ReservationResponse{
		ID:            reservation.ID,
		EquipmentID:   reservation.EquipmentID,
		EquipmentName: equipmentName,
		UserID:        reservation.UserID,
		UserName:      userName,
		StartTime:     reservation.StartTime,
		EndTime:       reservation.EndTime,
		Purpose:       reservation.Purpose,
		Status:        string(reservation.Status),
		ApproverID:    reservation.ApproverID,
		ApproverName:  approverName,
		CreatedAt:     reservation.CreatedAt,
	}
}

func mapDisposalItem(item service.DisposalItem) dto.DisposalResponse {
	request := item.Request
	equipmentName := ""
	if request.Equipment != nil {
		equipmentName = request.Equipment.Name
	}
	applicantName := ""
	if request.Applicant != nil {
		applicantName = request.Applicant.Name
	}
	approverName := ""
	if request.Approver != nil {
		approverName = request.Approver.Name
	}
	blockers := make([]dto.DisposalBlocker, 0, len(item.Blockers))
	for _, blocker := range item.Blockers {
		blockers = append(blockers, dto.DisposalBlocker{
			Type:     blocker.Type,
			RecordID: blocker.RecordID,
			Status:   blocker.Status,
			UserName: blocker.UserName,
			Detail:   blocker.Detail,
		})
	}
	return dto.DisposalResponse{
		ID:            request.ID,
		EquipmentID:   request.EquipmentID,
		EquipmentName: equipmentName,
		Reason:        request.Reason,
		Status:        string(request.Status),
		ApplicantID:   request.ApplicantID,
		ApplicantName: applicantName,
		ApproverID:    request.ApproverID,
		ApproverName:  approverName,
		Blockers:      blockers,
		ProcessedAt:   request.ProcessedAt,
		CreatedAt:     request.CreatedAt,
	}
}

func mapCategory(category model.EquipmentCategory) dto.CategoryResponse {
	children := make([]dto.CategoryResponse, 0, len(category.Children))
	for _, child := range category.Children {
		children = append(children, mapCategory(child))
	}
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
		Icon:        category.Icon,
		Children:    children,
	}
}

func mapAudit(log model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{
		ID:           log.ID,
		UserID:       log.UserID,
		UserName:     log.UserName,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Detail:       log.Detail,
		IP:           log.IP,
		CreatedAt:    log.CreatedAt,
	}
}
