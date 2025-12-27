package model

import (
	"errors"
	"project-manager/pkg"

	"gorm.io/gorm"
)

type Team struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Desc      *string        `json:"desc,omitempty"`
	Leader    *User          `json:"leader,omitempty" gorm:"foreignKey:LeaderID;references:ID"`
	LeaderID  *uint          `json:"leader_id,omitempty"`
	Projects  []*TeamProject `json:"projects,omitempty"` //invalid field found for struct project-manager/model.Team's field Projects: define a valid foreign key for relations or implement the Valuer/Scanner interface
	Users     []*User        `json:"users,omitempty" gorm:"many2many:team_users;"`
	CreatedAt int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64          `json:"updated_at" gorm:"autoUpdateTime"`
}

type TeamProject struct {
	ID        int      `json:"id" gorm:"primaryKey"`
	TeamID    uint     `json:"team_id" gorm:""`
	ProjectID uint     `json:"project_id" gorm:""`
	Name      string   `json:"name" gorm:""`
	Project   *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

type TeamModel struct{}

// Exist
func (m *TeamModel) GetByTeamName(name string) (*Team, error) {
	var dataObj Team
	err := pkg.DB.Debug().Where("name = ?", name).First(&dataObj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dataObj, err
}

// IsTeamLeader
func (m *TeamModel) IsTeamLeader(userID uint, teamID uint) (bool, error) {
	var team Team
	err := pkg.DB.Debug().Where("id = ? AND leader_id = ?", teamID, userID).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, err
}

// IsUserInTeam
func (m *TeamModel) IsUserInTeam(userID uint, teamID uint) (bool, error) {
	//查询team_users表
	var count int64
	err := pkg.DB.Debug().Table("team_users").Where("team_id = ? AND user_id = ?", teamID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByID
func (m *TeamModel) GetByID(id uint) (*Team, error) {
	var dataObj Team
	err := pkg.DB.Debug().Where("id = ?", id).First(&dataObj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dataObj, err
}

func (m *TeamModel) GetDetailByID(id uint) (*Team, error) {
	var dataObj Team
	err := pkg.DB.Debug().
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "team_id", "project_id", "name")
		}).
		Preload("Leader", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "nickname", "email", "logo", "created_at", "updated_at")
		}).
		Preload("Leader.Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "type", "desc")
		}).
		Where("id = ?", id).First(&dataObj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dataObj, err
}

// AddUserToTeam ...
func (m *TeamModel) AddUserToTeam(userID uint, teamID uint) error {
	//插入team_users表
	return pkg.DB.Debug().Table("team_users").Create(map[string]any{
		"team_id": teamID,
		"user_id": userID,
	}).Error
}

func (m *TeamModel) AddProjectToTeam(teamID uint, projectName string, projectDesc *string) error {
	newProject := &Project{
		Name: projectName,
		Desc: projectDesc,
	}
	err := pkg.DB.Create(newProject).Error
	if err != nil {
		return err
	}
	//关联project到team
	teamProject := &TeamProject{
		TeamID:    teamID,
		ProjectID: newProject.ID,
		Name:      newProject.Name,
	}
	err = pkg.DB.Create(teamProject).Error
	if err != nil {
		return err
	}
	return nil
}

// GetTeamsLedByUser 获取用户担任 Leader 的所有 Team
func (m *TeamModel) GetTeamsLedByUser(userID uint) ([]Team, error) {
	var teams []Team
	err := pkg.DB.Debug().Where("leader_id = ?", userID).Find(&teams).Error
	return teams, err
}

// GetUsersInTeams 获取指定 Teams 中的所有用户ID（用于可见性判断）
func (m *TeamModel) GetUsersInTeams(teamIDs []uint) ([]uint, error) {
	var userIDs []uint
	err := pkg.DB.Debug().Table("team_users").
		Where("team_id IN ?", teamIDs).
		Distinct("user_id").
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// list 查询条件 userID orderby page pagesize
func (m *TeamModel) List(userID uint, orderBy string, page int, pageSize int) ([]Team, error) {
	var teams []Team
	offset := (page - 1) * pageSize
	//查询列表 需要preload关联数据
	query := pkg.DB.Debug().Model(&Team{}).Preload("Leader").Preload("Leader.Roles")

	//如果userID不为0，表示查询该用户所在的team列表
	if userID != 0 {
		query = query.Joins("JOIN team_users ON team_users.team_id = teams.id").
			Where("team_users.user_id = ?", userID)
	}
	//排序
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	//分页
	if pageSize > 0 {
		query = query.Offset(offset).Limit(pageSize)
	}
	err := query.Find(&teams).Error
	return teams, err
}

// ListUsersByTeamID ...
func (m *TeamModel) ListUsersByTeamID(teamID uint, orderBy string, page int, pageSize int, name string) ([]User, error) {
	var users []User
	offset := (page - 1) * pageSize
	//查询列表 需要preload关联数据
	query := pkg.DB.Debug().Model(&User{}).Joins("JOIN team_users ON team_users.user_id = users.id").
		Where("team_users.team_id = ?", teamID).
		Preload("Roles")
	if name != "" {
		query = query.Where("users.username LIKE ? OR users.nickname LIKE ?", "%"+name+"%", "%"+name+"%")
	}
	//排序
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	//分页
	if pageSize > 0 {
		query = query.Offset(offset).Limit(pageSize)
	}
	err := query.Find(&users).Error
	return users, err
}

// DeleteTeam ...
func (m *TeamModel) DeleteTeam(teamID uint) error {
	//删除team_projects关联数据
	err := pkg.DB.Debug().Where("team_id = ?", teamID).Delete(&TeamProject{}).Error
	if err != nil {
		return err
	}
	//删除team_users关联数据
	err = pkg.DB.Debug().Table("team_users").Where("team_id = ?", teamID).Delete(nil).Error
	if err != nil {
		return err
	}
	//删除team数据
	err = pkg.DB.Debug().Where("id = ?", teamID).Delete(&Team{}).Error
	return err
}
