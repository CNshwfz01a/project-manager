package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     *string   `gorm:"size:100" json:"email,omitempty"`
	Nickname  *string   `gorm:"size:50" json:"nickname,omitempty"`
	Logo      *string   `gorm:"size:255" json:"logo,omitempty"`
	Roles     []*Role   `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


