package response

import "project-manager/model"

type RoleListRsp struct {
	Roles []model.Role `json:"list"`
	Total int64        `json:"total"`
}

// type RoleAddRsp struct {
// 	Role model.Role `json:"role"`
// }
