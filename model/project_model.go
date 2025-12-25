package model

type Project struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	Name      string  `json:"name"`
	Desc      *string `json:"desc" gorm:"size:255;default:''"`
	Status    string  `json:"status" gorm:"size:20;not null"` // WAIT_FOR_SCHEDULE, IN_PROGRESS, FINISHED
	Users     []*User `json:"users,omitempty" gorm:"many2many:project_users;"`
	CreatedAt int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64   `json:"updated_at" gorm:"autoUpdateTime"`
}

type ProjectModel struct{}
