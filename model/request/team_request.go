package request

type TeamAddReq struct {
	Name string  `json:"name" validate:"required"`
	Desc *string `json:"desc,omitempty"`
}

type TeamAddUserReq struct {
	TeamID uint `json:"team_id"`
	UserID uint `json:"user_id" validate:"required"`
}

type TeamAddProjectReq struct {
	TeamID      uint    `json:"team_id"`
	ProjectName string  `json:"name" validate:"required"`
	ProjectDesc *string `json:"desc,omitempty"`
}

/*
*
// TeamPatch 示例: [{"op":"replace","path":"/leader","value":{"id":5}}]
*/
type TeamPatch []struct {
	Op    string      `json:"op" validate:"required,oneof=replace add remove"`
	Path  string      `json:"path" validate:"required,oneof=/name /desc /leader"`
	Value interface{} `json:"value,omitempty"`
}

type TeamUpdateReq struct {
	Name *string `json:"name,omitempty"`
	Desc *string `json:"desc,omitempty"`
}
