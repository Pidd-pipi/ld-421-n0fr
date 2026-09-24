package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/labequipment/lab-equipment/database/migrations"
	"github.com/labequipment/lab-equipment/database/seeds"
	"github.com/labequipment/lab-equipment/internal/config"
	"github.com/labequipment/lab-equipment/internal/handler"
	"github.com/labequipment/lab-equipment/internal/repository"
	"github.com/labequipment/lab-equipment/internal/service"
	"gorm.io/gorm"
)

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func doReq(t *testing.T, srv *httptest.Server, method, path, token string, body any) (int, envelope) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, srv.URL+path, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("decode envelope: %v, raw: %s", err, raw)
	}
	return resp.StatusCode, env
}

func TestDisposalHTTPFlow(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := migrations.AutoMigrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := seeds.Seed(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	log := slog.Default()

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	equipmentRepo := repository.NewEquipmentRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)
	maintenanceRepo := repository.NewMaintenanceRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	disposalRepo := repository.NewDisposalRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	auditService := service.NewAuditService(auditRepo, log)
	authService := service.NewAuthService(userRepo, config.JWTConfig{Secret: "test-secret", ExpireHours: 1}, log)
	userService := service.NewUserService(userRepo, log)
	categoryService := service.NewCategoryService(categoryRepo, log)
	equipmentService := service.NewEquipmentService(equipmentRepo, categoryRepo, userRepo, auditService, log)
	borrowService := service.NewBorrowService(borrowRepo, equipmentRepo, disposalRepo, auditService, log)
	maintenanceService := service.NewMaintenanceService(maintenanceRepo, equipmentRepo, auditService, log)
	reservationService := service.NewReservationService(reservationRepo, equipmentRepo, disposalRepo, auditService, log)
	disposalService := service.NewDisposalService(disposalRepo, equipmentRepo, borrowRepo, reservationRepo, auditService, log)
	dashboardService := service.NewDashboardService(equipmentRepo, borrowRepo, reservationRepo, log)

	engine := NewRouter(Dependencies{
		AuthHandler:        handler.NewAuthHandler(authService),
		UserHandler:        handler.NewUserHandler(userService),
		CategoryHandler:    handler.NewCategoryHandler(categoryService),
		EquipmentHandler:   handler.NewEquipmentHandler(equipmentService),
		BorrowHandler:      handler.NewBorrowHandler(borrowService),
		MaintenanceHandler: handler.NewMaintenanceHandler(maintenanceService),
		ReservationHandler: handler.NewReservationHandler(reservationService),
		DisposalHandler:    handler.NewDisposalHandler(disposalService),
		DashboardHandler:   handler.NewDashboardHandler(dashboardService),
		AuditHandler:       handler.NewAuditHandler(auditService),
		AuthService:        authService,
		AuditService:       auditService,
		Logger:             log,
	})
	srv := httptest.NewServer(engine)
	defer srv.Close()

	// 登录（admin / admin123）
	status, env := doReq(t, srv, http.MethodPost, "/api/v1/auth/login", "", map[string]any{"username": "admin", "password": "admin123"})
	if status != 200 {
		t.Fatalf("login status %d: %s", status, env.Message)
	}
	var loginData struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(env.Data, &loginData); err != nil || loginData.Token == "" {
		t.Fatalf("login token missing: %s", env.Data)
	}
	token := loginData.Token

	// 未认证访问被拒
	if status, _ := doReq(t, srv, http.MethodGet, "/api/v1/equipment/1/disposals", "", nil); status != 401 {
		t.Fatalf("expected 401 without token, got %d", status)
	}

	// 学生角色不能提交报废申请
	_, studentEnv := doReq(t, srv, http.MethodPost, "/api/v1/auth/login", "", map[string]any{"username": "student", "password": "stu123"})
	var studentLogin struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(studentEnv.Data, &studentLogin)
	if status, _ := doReq(t, srv, http.MethodPost, "/api/v1/equipment/1/disposals", studentLogin.Token, map[string]any{"reason": "x"}); status != 403 {
		t.Fatalf("expected 403 for student submit, got %d", status)
	}

	// 借用设备 1 并审批通过，制造未归还阻塞项
	status, env = doReq(t, srv, http.MethodPost, "/api/v1/borrows", token, map[string]any{
		"equipmentId": 1, "borrowDate": "2026-09-24T10:00:00Z", "expectedReturnDate": "2026-10-01T10:00:00Z", "reason": "实验",
	})
	if status != 200 {
		t.Fatalf("create borrow status %d: %s", status, env.Message)
	}
	var borrow struct {
		ID uint `json:"id"`
	}
	_ = json.Unmarshal(env.Data, &borrow)
	if status, env = doReq(t, srv, http.MethodPost, fmt.Sprintf("/api/v1/borrows/%d/approve", borrow.ID), token, nil); status != 200 {
		t.Fatalf("approve borrow status %d: %s", status, env.Message)
	}

	// 提交报废申请
	status, env = doReq(t, srv, http.MethodPost, "/api/v1/equipment/1/disposals", token, map[string]any{"reason": "设备老化"})
	if status != 200 {
		t.Fatalf("submit disposal status %d: %s", status, env.Message)
	}
	var disposal struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}
	_ = json.Unmarshal(env.Data, &disposal)
	if disposal.Status != "Pending" {
		t.Fatalf("expected Pending, got %s", disposal.Status)
	}

	// 重复提交 → 409
	if status, _ = doReq(t, srv, http.MethodPost, "/api/v1/equipment/1/disposals", token, map[string]any{"reason": "重复"}); status != 409 {
		t.Fatalf("expected 409 on duplicate submit, got %d", status)
	}

	// 新借用被暂停 → 409
	if status, _ = doReq(t, srv, http.MethodPost, "/api/v1/borrows", token, map[string]any{
		"equipmentId": 1, "borrowDate": "2026-09-24T10:00:00Z", "expectedReturnDate": "2026-10-01T10:00:00Z",
	}); status != 409 {
		t.Fatalf("expected 409 on borrow during disposal, got %d", status)
	}

	// 新预约被暂停 → 409
	if status, _ = doReq(t, srv, http.MethodPost, "/api/v1/reservations", token, map[string]any{
		"equipmentId": 1, "startTime": "2026-09-25T10:00:00Z", "endTime": "2026-09-25T12:00:00Z", "purpose": "实验",
	}); status != 409 {
		t.Fatalf("expected 409 on reservation during disposal, got %d", status)
	}

	// 列表展示进度与实时阻塞项
	status, env = doReq(t, srv, http.MethodGet, "/api/v1/equipment/1/disposals", token, nil)
	if status != 200 {
		t.Fatalf("list disposals status %d: %s", status, env.Message)
	}
	var items []struct {
		ID       uint `json:"id"`
		Status   string
		Reason   string
		Blockers []struct {
			Type string `json:"type"`
		}
	}
	_ = json.Unmarshal(env.Data, &items)
	if len(items) != 1 || items[0].Status != "Pending" || items[0].Reason != "设备老化" {
		t.Fatalf("unexpected disposal list: %s", env.Data)
	}
	if len(items[0].Blockers) != 1 || items[0].Blockers[0].Type != "borrow" {
		t.Fatalf("expected live borrow blocker, got %s", env.Data)
	}

	// 审批：有未归还借用 → 退回并列出阻塞项
	status, env = doReq(t, srv, http.MethodPost, fmt.Sprintf("/api/v1/equipment/1/disposals/%d/approve", disposal.ID), token, nil)
	if status != 200 {
		t.Fatalf("approve disposal status %d: %s", status, env.Message)
	}
	var returned struct {
		Status   string
		Blockers []struct {
			Type string `json:"type"`
		}
	}
	_ = json.Unmarshal(env.Data, &returned)
	if returned.Status != "Rejected" || len(returned.Blockers) != 1 {
		t.Fatalf("expected returned with blockers, got %s", env.Data)
	}

	// 重复审批 → 409
	if status, _ = doReq(t, srv, http.MethodPost, fmt.Sprintf("/api/v1/equipment/1/disposals/%d/approve", disposal.ID), token, nil); status != 409 {
		t.Fatalf("expected 409 on repeated approve, got %d", status)
	}

	// 设备仍未报废
	status, env = doReq(t, srv, http.MethodGet, "/api/v1/equipment/1", token, nil)
	var eq struct {
		Status string `json:"status"`
	}
	_ = json.Unmarshal(env.Data, &eq)
	if eq.Status == "Retired" {
		t.Fatalf("equipment must not be Retired while blocked")
	}

	// 归还借用 → 重新提交 → 审批通过 → 设备 Retired
	if status, env = doReq(t, srv, http.MethodPost, fmt.Sprintf("/api/v1/borrows/%d/return", borrow.ID), token, map[string]any{
		"actualReturnDate": "2026-09-24T18:00:00Z", "returnCondition": "Good",
	}); status != 200 {
		t.Fatalf("return borrow status %d: %s", status, env.Message)
	}
	status, env = doReq(t, srv, http.MethodPost, "/api/v1/equipment/1/disposals", token, map[string]any{"reason": "设备老化，再次申请"})
	if status != 200 {
		t.Fatalf("resubmit disposal status %d: %s", status, env.Message)
	}
	_ = json.Unmarshal(env.Data, &disposal)
	status, env = doReq(t, srv, http.MethodPost, fmt.Sprintf("/api/v1/equipment/1/disposals/%d/approve", disposal.ID), token, nil)
	if status != 200 {
		t.Fatalf("approve disposal status %d: %s", status, env.Message)
	}
	_ = json.Unmarshal(env.Data, &returned)
	if returned.Status != "Approved" {
		t.Fatalf("expected Approved, got %s", env.Data)
	}
	status, env = doReq(t, srv, http.MethodGet, "/api/v1/equipment/1", token, nil)
	_ = json.Unmarshal(env.Data, &eq)
	if eq.Status != "Retired" {
		t.Fatalf("expected Retired, got %s", eq.Status)
	}

	// 原借用记录保留
	status, env = doReq(t, srv, http.MethodGet, fmt.Sprintf("/api/v1/borrows/%d", borrow.ID), token, nil)
	var kept struct {
		Status string `json:"status"`
	}
	_ = json.Unmarshal(env.Data, &kept)
	if kept.Status != "Returned" {
		t.Fatalf("borrow record should be kept as Returned, got %s", env.Data)
	}

	// 已报废设备再次提交 → 409
	if status, _ = doReq(t, srv, http.MethodPost, "/api/v1/equipment/1/disposals", token, map[string]any{"reason": "再次"}); status != 409 {
		t.Fatalf("expected 409 on submit for retired equipment, got %d", status)
	}
}
