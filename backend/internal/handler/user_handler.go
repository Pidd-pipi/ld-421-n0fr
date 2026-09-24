package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/repository"
	"github.com/labequipment/lab-equipment/internal/service"
)

// UserHandler 用户处理器。
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 构造用户处理器。
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List 分页查询用户。
func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	roleID, _ := strconv.ParseUint(c.Query("role_id"), 10, 64)
	filter := repository.UserFilter{
		Keyword:    c.Query("keyword"),
		RoleID:     uint(roleID),
		Pagination: repository.Pagination{Page: page, PageSize: pageSize},
	}
	users, total, err := h.userService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, mapUser(user))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}
