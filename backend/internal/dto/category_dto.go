package dto

// CreateCategoryRequest 新增分类。
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description" binding:"omitempty,max=512"`
	Icon        string `json:"icon" binding:"omitempty,max=128"`
}

// UpdateCategoryRequest 更新分类。
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description" binding:"omitempty,max=512"`
	Icon        string `json:"icon" binding:"omitempty,max=128"`
}

// CategoryResponse 分类返回。
type CategoryResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	ParentID    *uint              `json:"parentId"`
	Description string             `json:"description"`
	Icon        string             `json:"icon"`
	Children    []CategoryResponse `json:"children,omitempty"`
}
