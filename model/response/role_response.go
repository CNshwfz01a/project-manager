package response

import "project-manager/model"

type RoleListRsp struct {
	Total int64        `json:"total"`
	Roles []model.Role `json:"roles"`
}
