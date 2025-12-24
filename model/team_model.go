package model

import (
	"errors"
	"project-manager/pkg"

	"gorm.io/gorm"
)

type Team struct {
	ID        int            `json:"id" gorm:"primaryKey"`
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
	TeamID    int    `json:"team_id" gorm:""`
	ProjectID int    `json:"project_id" gorm:""`
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
func (m *TeamModel) IsTeamLeader(userID uint, teamID int) (bool, error) {
	var team Team
	err := pkg.DB.Debug().Where("id = ? AND leader_id = ?", teamID, userID).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, err
}


