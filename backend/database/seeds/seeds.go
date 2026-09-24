package seeds

import (
	"fmt"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed 初始化角色、用户、分类与示例设备。
func Seed(db *gorm.DB) error {
	if err := seedRoles(db); err != nil {
		return err
	}
	if err := seedUsers(db); err != nil {
		return err
	}
	if err := seedCategories(db); err != nil {
		return err
	}
	if err := seedEquipment(db); err != nil {
		return err
	}
	return nil
}

func seedRoles(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Role{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count roles: %w", err)
	}
	if count > 0 {
		return nil
	}
	roles := []model.Role{
		{Code: "Admin", Name: "管理员", Description: "拥有全部权限"},
		{Code: "LabManager", Name: "实验室管理员", Description: "负责审批与设备管理"},
		{Code: "Researcher", Name: "研究员", Description: "可借用与维护设备"},
		{Code: "Student", Name: "学生", Description: "可借用与预约设备"},
	}
	return db.Create(&roles).Error
}

func seedUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}
	var roles []model.Role
	if err := db.Find(&roles).Error; err != nil {
		return err
	}
	roleByCode := make(map[string]uint, len(roles))
	for _, role := range roles {
		roleByCode[role.Code] = role.ID
	}
	users := []struct {
		username string
		name     string
		password string
		role     string
		email    string
	}{
		{username: "admin", name: "系统管理员", password: "admin123", role: "Admin", email: "admin@lab.local"},
		{username: "labmanager", name: "实验室管理员", password: "lab123", role: "LabManager", email: "manager@lab.local"},
		{username: "researcher", name: "研究员", password: "res123", role: "Researcher", email: "researcher@lab.local"},
		{username: "student", name: "学生", password: "stu123", role: "Student", email: "student@lab.local"},
	}
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", u.username, err)
		}
		user := model.User{
			Username:     u.username,
			PasswordHash: string(hash),
			Name:         u.name,
			Email:        u.email,
			RoleID:       roleByCode[u.role],
			Active:       true,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("seed user %s: %w", u.username, err)
		}
	}
	return nil
}

func seedCategories(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.EquipmentCategory{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count categories: %w", err)
	}
	if count > 0 {
		return nil
	}
	root := model.EquipmentCategory{Name: "分析仪器", Description: "用于样品分析的仪器", Icon: "ExperimentOutlined"}
	if err := db.Create(&root).Error; err != nil {
		return err
	}
	children := []model.EquipmentCategory{
		{Name: "光谱仪", ParentID: &root.ID, Description: "光谱分析设备", Icon: "FundProjectionScreenOutlined"},
		{Name: "色谱仪", ParentID: &root.ID, Description: "色谱分析设备", Icon: "DotChartOutlined"},
		{Name: "显微镜", ParentID: &root.ID, Description: "显微观察设备", Icon: "EyeOutlined"},
	}
	return db.Create(&children).Error
}

func seedEquipment(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Equipment{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count equipment: %w", err)
	}
	if count > 0 {
		return nil
	}
	var categories []model.EquipmentCategory
	if err := db.Find(&categories).Error; err != nil {
		return err
	}
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	categoryIDByName := make(map[string]uint, len(categories))
	for _, category := range categories {
		categoryIDByName[category.Name] = category.ID
	}
	categoryID := func(name string) uint {
		if id, ok := categoryIDByName[name]; ok {
			return id
		}
		return 0
	}
	ownerID := users[0].ID
	for _, u := range users {
		if u.Username == "researcher" {
			ownerID = u.ID
		}
	}
	now := time.Now()
	in15Days := now.AddDate(0, 0, 15)
	equipment := []model.Equipment{
		{Name: "紫外可见分光光度计", Code: "EQ-001", CategoryID: categoryID("光谱仪"), BrandModel: "Shimadzu UV-2600", SerialNumber: "UV2600-001", PurchaseDate: &now, PurchasePrice: 128000, Location: "A栋-302-01", Status: constants.AssetStatusAvailable, OwnerID: ownerID, Supplier: "岛津公司", WarrantyExpiry: &in15Days, ImageURL: ""},
		{Name: "气相色谱仪", Code: "EQ-002", CategoryID: categoryID("色谱仪"), BrandModel: "Agilent 7890B", SerialNumber: "7890B-002", PurchasePrice: 350000, Location: "A栋-302-02", Status: constants.AssetStatusAvailable, OwnerID: ownerID, Supplier: "安捷伦", WarrantyExpiry: &in15Days},
		{Name: "倒置荧光显微镜", Code: "EQ-003", CategoryID: categoryID("显微镜"), BrandModel: "Olympus IX73", SerialNumber: "IX73-003", PurchasePrice: 210000, Location: "B栋-105-01", Status: constants.AssetStatusMaintenance, OwnerID: ownerID, Supplier: "奥林巴斯"},
	}
	return db.Create(&equipment).Error
}
