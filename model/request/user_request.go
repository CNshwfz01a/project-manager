package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserLoginReq struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required"`
}

type UserChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password"` //8-30个字符，正则表达式只能包含字母、数字、下划线和连字符
}

type UserAddReq struct {
	Username string `json:"username" validate:"required,min=4,max=30,password"` //4-30个字符，正则表达式只能包含字母、数字、下划线和连字符
	Password string `json:"password" validate:"required,min=8,max=30,password"` //8-30个字符，正则表达式只能包含字母、数字、下划线和连字符
}

type UserListReq struct {
	//排序字段enum "created_at" "updated_at"
	OrderBy string `json:"order_by" form:"order_by" validate:"omitempty,oneof=created_at updated_at"`
	//page 不传则默认1
	Page int `json:"page" form:"page" validate:"omitempty,min=1"`
	//page size 不传则默认10
	PageSize int `json:"page_size" form:"page_size" validate:"omitempty,min=1,max=100"`
	//搜索关键词name 包含username和nickname
	Name string `json:"name" form:"name" validate:"omitempty"`
	//数组team_id
	TeamIDs []int `json:"team_id" form:"team_id[]" validate:"omitempty,dive,gt=0"`
	//数组role_name
	RoleNames []string `json:"role_name" form:"role_name[]" validate:"omitempty,dive,required"`
}

// 设置默认值
// SetDefaults 设置默认值
func (r *UserListReq) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
}

// SetParams 设置查询参数
func (r *UserListReq) SetParams(c *gin.Context) {
	r.OrderBy = c.DefaultQuery("order_by", r.OrderBy)
	r.Page, _ = strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(r.Page)))
	r.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(r.PageSize)))
	r.Name = c.DefaultQuery("name", r.Name)
	//如果有TeamIDs参数
	teamIDStrs, _ := c.GetQueryArray("team_id[]") //如果没有该参数则返回空切片
	r.TeamIDs = make([]int, 0, len(teamIDStrs))
	for _, idStr := range teamIDStrs {
		if id, err := strconv.Atoi(idStr); err == nil {
			r.TeamIDs = append(r.TeamIDs, id)
		}
	}
	//如果有RoleNames参数
	r.RoleNames, _ = c.GetQueryArray("role_name[]")
}

type UserAssignRoleReq struct {
	RoleID uint `json:"role_id" validate:"required,gt=0"`
}

type UserUpdateProfileReq struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Logo     string `json:"logo" validate:"omitempty,url,max=255"`
}
