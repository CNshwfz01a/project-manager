package model

import (
	"project-manager/model/request"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	ID        uint    `gorm:"primarykey" json:"id"`
	Username  string  `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Status    uint8   `gorm:"not null;default:0" json:"status"` // 0:未激活 1:激活
	Password  string  `gorm:"size:255;not null" json:"-"`
	Email     *string `gorm:"size:100" json:"email,omitempty"`
	Nickname  *string `gorm:"size:50" json:"nickname,omitempty"`
	Logo      *string `gorm:"size:255" json:"logo,omitempty"`
	Roles     []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64   `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserModel struct{}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	err := pkg.DB.Preload("Roles").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (m *UserModel) GetByUsernameOrEmail(loginID string) (*User, error) {
	var user User
	err := pkg.DB.Preload("Roles").Where("username = ? OR email = ?", loginID, loginID).First(&user).Error
	return &user, err
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := pkg.DB.Preload("Roles").Where("email = ?", email).First(&user).Error
	return &user, err
}

// 查询用户详情
func (m *UserModel) GetByID(id uint) (*User, error) {
	var user User
	err := pkg.DB.Preload("Roles").First(&user, id).Error
	return &user, err
}

// DeleteRelationByRoleID 删除用户与角色的关联关系
func (m *UserModel) DeleteRelationByRoleID(roleID uint) error {
	return pkg.DB.Table("user_roles").Where("role_id = ?", roleID).Delete(nil).Error
}

/*
List 查询用户列表
查询参数见 request.UserListReq
用户角色限定显示范围 1.admin角色可以查看所有用户 2.非admin角色只能查看自己相同team下的用户
*/
func (m *UserModel) List(c *gin.Context, req *request.UserListReq) ([]User, error) {
	var users []User
	//构建查询
	db := pkg.DB.Model(&User{}).Preload("Roles").Select("id", "username", "status", "email", "nickname", "logo", "created_at", "updated_at")

	//使用role_service的GetRoleByName函数判断当前用户角色
	isAdmin, roleErr := GetRoleByName("admin", c)
	if roleErr != nil {
		return nil, roleErr.(error)
	}

	if !isAdmin {
		// 非admin角色 只能查看相同team下的用户
		userID, exists := c.Get("user_id")
		if !exists {
			return nil, pkg.NewRspError(401, nil)
		}
		var teamIDs []uint
		// 查询该用户所属的团队ID列表
		err := pkg.DB.Table("team_users").Select("team_id").Where("user_id = ?", userID.(uint)).Scan(&teamIDs).Error
		if err != nil {
			return nil, err
		}
		//如果没有所属团队，则返回自己
		if len(teamIDs) == 0 {
			db = db.Where("id = ?", userID.(uint))
		} else {
			db = db.Joins("JOIN team_users ON team_users.user_id = users.id").
				Where("team_users.team_id IN ?", teamIDs)
		}
	}
	// admin角色可以查看所有用户（不添加额外的where条件）

	// 按名称搜索（username或nickname）
	if req.Name != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}

	// 按team_id筛选（如果User表有team_id字段）
	if len(req.TeamIDs) > 0 {
		db = db.Joins("JOIN team_users ON team_users.user_id = users.id").
			Where("team_users.team_id IN ?", req.TeamIDs)
	}

	// 按角色名称筛选
	if len(req.RoleNames) > 0 {
		db = db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.name IN ?", req.RoleNames)
	}

	if req.OrderBy != "" {
		db = db.Order(req.OrderBy + " DESC")
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	db = db.Offset(offset).Limit(req.PageSize)

	// 执行查询
	err := db.Find(&users).Error

	return users, err
}

// DeleteByID 根据ID删除用户
func (m *UserModel) DeleteByID(id uint) error {
	//删除用户角色关联 team关联 project关联等
	err := pkg.DB.Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Table("user_roles").Where("user_id = ?", id).Delete(nil).Error; err != nil {
			return err
		}
		// 删除用户团队关联
		if err := tx.Table("team_users").Where("user_id = ?", id).Delete(nil).Error; err != nil {
			return err
		}
		// 删除用户项目关联
		if err := tx.Table("project_users").Where("user_id = ?", id).Delete(nil).Error; err != nil {
			return err
		}
		// 如果用户是Team Leader，将Team的LeaderID置为NULL
		if err := tx.Model(&Team{}).Where("leader_id = ?", id).Update("leader_id", nil).Error; err != nil {
			return err
		}
		// 删除用户
		if err := tx.Delete(&User{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
