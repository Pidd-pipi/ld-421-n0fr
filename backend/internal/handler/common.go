package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/labequipment/lab-equipment/internal/constants"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/repository"
)

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{Code: constants.CodeSuccess, Message: "ok", Data: data})
}

func fail(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, response{Code: code, Message: message, Data: nil})
}

// failError 将错误转换为统一响应。
func failError(c *gin.Context, err error) {
	var be *apperrors.BusinessError
	if errors.As(err, &be) {
		fail(c, be.HTTPStatus, be.Code, be.Message)
		return
	}
	var ve *apperrors.ValidationError
	if errors.As(err, &ve) {
		c.JSON(http.StatusUnprocessableEntity, response{Code: constants.CodeValidation, Message: ve.Message, Data: ve.FieldErrors})
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		fail(c, http.StatusNotFound, constants.CodeNotFound, "资源不存在")
		return
	}
	fail(c, http.StatusInternalServerError, constants.CodeInternal, "服务器内部错误")
}

// bindJSON 绑定 JSON 并将校验错误转换为统一响应。
func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			fields := make(map[string]string, len(verr))
			for _, fe := range verr {
				fields[fe.Field()] = fe.Tag()
			}
			fail(c, http.StatusUnprocessableEntity, constants.CodeValidation, "参数校验失败")
			return false
		}
		fail(c, http.StatusBadRequest, constants.CodeBadRequest, "请求参数格式错误")
		return false
	}
	return true
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		fail(c, http.StatusBadRequest, constants.CodeBadRequest, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

func parsePage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	return page, pageSize
}
