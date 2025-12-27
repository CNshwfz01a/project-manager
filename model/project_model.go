package model

import (
	"errors"
	"project-manager/pkg"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	Name      string  `json:"name" gorm:"size:100;not null"`
	Desc      *string `json:"desc" gorm:"size:255"`
	Status    string  `json:"status" gorm:"size:20;not null;default:'WAIT_FOR_SCHEDULE'"` // WAIT_FOR_SCHEDULE, IN_PROGRESS, FINISHED
	Users     []*User `json:"users,omitempty" gorm:"many2many:project_users"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

type ProjectModel struct{}

// GetByID 获取项目详情
func (m *ProjectModel) GetByID(id uint) (*Project, error) {
	var project Project
	err := pkg.DB.Debug().Where("id = ?", id).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &project, err
}

// GetByName 根据名称获取项目
func (m *ProjectModel) GetByName(name string) (*Project, error) {
	var project Project
	err := pkg.DB.Debug().Where("name = ?", name).First(&project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &project, err
}

// GetTeamIDByProjectID 根据项目ID获取所属的TeamID
func (m *ProjectModel) GetTeamIDByProjectID(projectID uint) (uint, error) {
	var teamProject TeamProject
	err := pkg.DB.Debug().Where("project_id = ?", projectID).First(&teamProject).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return teamProject.TeamID, err
}

// IsUserInProject 检查用户是否在项目中
func (m *ProjectModel) IsUserInProject(userID uint, projectID uint) (bool, error) {
	var count int64
	err := pkg.DB.Debug().Table("project_users").Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AddUserToProject 添加用户到项目
func (m *ProjectModel) AddUserToProject(userID uint, projectID uint) error {
	return pkg.DB.Debug().Table("project_users").Create(map[string]any{
		"project_id": projectID,
		"user_id":    userID,
	}).Error
}
