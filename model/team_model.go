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
	ID        int    `json:"id" gorm:"primaryKey"`
	TeamID    uint   `json:"team_id" gorm:""`
	ProjectID uint   `json:"project_id" gorm:""`
	Name      string `json:"name" gorm:""`
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
