package response

import (
	"project-manager/model/request"
)

type UserListResp struct {
	ID        uint                  `json:"id" form:"id"`
	Username  string                `json:"username" form:"username"`
	Status    uint8                 `json:"status" form:"status"`
	Email     *string               `json:"email" form:"email"`
	Nickname  *string               `json:"nickname" form:"nickname"`
	Logo      *string               `json:"logo" form:"logo"`
	Roles     []request.RoleListReq `json:"roles" form:"roles"`
	CreatedAt int64                 `json:"created_at" form:"created_at"`
	UpdatedAt int64                 `json:"updated_at" form:"updated_at"`
}
