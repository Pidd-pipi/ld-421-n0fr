package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/service"
)

// CategoryHandler 分类处理器。
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 构造分类处理器。
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// List 查询分类树。
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.ListTree(c.Request.Context())
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		items = append(items, mapCategory(category))
	}
	ok(c, items)
}

// Create 新增分类。
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if !bindJSON(c, &req) {
		return
	}
	category := &model.EquipmentCategory{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		Icon:        req.Icon,
	}
	created, err := h.categoryService.Create(c.Request.Context(), category)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapCategory(*created))
}

// Update 更新分类。
func (h *CategoryHandler) Update(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.UpdateCategoryRequest
	if !bindJSON(c, &req) {
		return
	}
	category := &model.EquipmentCategory{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		Icon:        req.Icon,
	}
	updated, err := h.categoryService.Update(c.Request.Context(), id, category)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapCategory(*updated))
}

// Delete 删除分类。
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}
