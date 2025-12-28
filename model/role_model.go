package model

import (
	"errors"
	"fmt"

	"project-manager/model/request"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Role struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Type string `gorm:"size:20;not null" json:"type"` // System or Custom
	Desc string `gorm:"size:255;default:''" json:"desc,omitempty"`
}

type RoleModel struct{}

// GetByID
func (m *RoleModel) GetByID(id uint) (*Role, error) {
	var role Role
	err := pkg.DB.First(&role, id).Error
	return &role, err
}

// Exist
func (m *RoleModel) Exist(filter map[string]any) bool {
	var dataObj Role
	err := pkg.DB.Debug().Order("id DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Add
func (m *RoleModel) Add(role *Role) error {
	return pkg.DB.Create(role).Error
}

// List
func (m *RoleModel) List(req *request.RoleListReq) ([]Role, error) {
	var list []Role
	//select id,name,type,desc
	db := pkg.DB.Model(&Role{}).Select("id", "name", "type", "desc").Order("id DESC")
	return list, db.Find(&list).Error
}

// Count
func (m *RoleModel) Count() (int64, error) {
	var count int64
	err := pkg.DB.Model(&Role{}).Count(&count).Error
	return count, err
}

// Delete
func (m *RoleModel) Delete(id uint) error {
	//先查询是否存在
	var r Role
	err := pkg.DB.First(&r, id).Error
	if err != nil {
		return err
	}
	return pkg.DB.Delete(&Role{}, id).Error
}

// GetByName
func (m *RoleModel) GetByName(name string) (*Role, error) {
	var role Role
	err := pkg.DB.Where("name = ?", name).First(&role).Error
	return &role, err
}

func GetRoleByName(name string, c *gin.Context) (isExist bool, repError any) {
	userRoleInterface, exists := c.Get("user_role")
	if !exists {
		return false, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userRoles := userRoleInterface.([]*Role)
	for _, role := range userRoles {
		if role.Name == name {
			isExist = true
			break
		}
	}
	return isExist, nil
}
