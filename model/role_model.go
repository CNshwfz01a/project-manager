package model

import (
	"project-manager/model/request"
	"project-manager/pkg"
	"time"
)

type Role struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Type      string    `gorm:"size:20;not null" json:"type"` // System or Custom
	Desc      *string   `gorm:"size:255" json:"desc,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleModel struct{}

// List
func (m *RoleModel) List(req *request.RoleListReq) ([]Role, error) {
	var list []Role
	db := pkg.DB.Model(&Role{})
	return list, db.Find(&list).Error
}

// Count
func (m *RoleModel) Count() (int64, error) {
	var count int64
	err := pkg.DB.Model(&Role{}).Count(&count).Error
	return count, err
}
