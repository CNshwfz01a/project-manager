package request

type ProjectAddUserReq struct {
	ProjectID uint `json:"project_id" validate:"required"`
	UserID    uint `json:"user_id" validate:"required"`
}
