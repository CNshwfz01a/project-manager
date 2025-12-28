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

// DeleteProject 删除项目
func (m *ProjectModel) DeleteProject(projectID uint) error {
	//先删除关联的用户
	err := pkg.DB.Debug().Table("project_users").Where("project_id = ?", projectID).Delete(nil).Error
	if err != nil {
		return err
	}
	//再删除关联的团队关系
	err = pkg.DB.Debug().Table("team_projects").Where("project_id = ?", projectID).Delete(nil).Error
	if err != nil {
		return err
	}
	//再删除项目
	return pkg.DB.Debug().Where("id = ?", projectID).Delete(&Project{}).Error
}

// RemoveUserFromProject 清退项目中的用户
func (m *ProjectModel) RemoveUserFromProject(userID uint, projectID uint) error {
	return pkg.DB.Debug().Table("project_users").Where("project_id = ? AND user_id = ?", projectID, userID).Delete(nil).Error
}

// GetUsersInProject 获取项目中的所有用户
func (m *ProjectModel) GetUsersInProject(projectID uint, orderBy string, page int, pageSize int, name string) ([]*User, error) {
	var users []*User
	query := pkg.DB.Debug().Table("users").
		Joins("JOIN project_users ON users.id = project_users.user_id").
		Where("project_users.project_id = ?", projectID).
		Preload("Roles")

	if name != "" {
		query = query.Where("users.username LIKE ?", "%"+name+"%")
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetProjectsByUserID 获取用户参与的所有项目
func (m *ProjectModel) GetProjectsByUserID(userID uint) ([]Project, error) {
	var projects []Project
	err := pkg.DB.Debug().Model(&Project{}).
		Joins("JOIN project_users ON project_users.project_id = projects.id").
		Where("project_users.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}

// GetCommonProjects 获取两个用户共同参与的项目
func (m *ProjectModel) GetCommonProjects(userID1, userID2 uint) ([]Project, error) {
	var projects []Project
	err := pkg.DB.Debug().Model(&Project{}).
		Joins("JOIN project_users as pu1 ON pu1.project_id = projects.id").
		Joins("JOIN project_users as pu2 ON pu2.project_id = projects.id").
		Where("pu1.user_id = ? AND pu2.user_id = ?", userID1, userID2).
		Find(&projects).Error
	return projects, err
}
