package response

import (
	"project-manager/model"
)

type TeamGetResp struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	Desc      string               `json:"desc"`
	Leader    *model.User          `json:"leader"`
	Projects  []*model.TeamProject `json:"projects"`
	CreatedAt int64                `json:"created_at"`
	UpdatedAt int64                `json:"updated_at"`
}

type TeamListResp struct {
	Teams []TeamGetResp `json:"teams"`
}
