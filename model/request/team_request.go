package request

type TeamAddReq struct {
	Name string  `json:"name" validate:"required"`
	Desc *string `json:"desc,omitempty"`
}

type TeamAddUserReq struct {
	TeamID int `json:"team_id" validate:"required"`
	UserID int `json:"user_id" validate:"required"`
}
