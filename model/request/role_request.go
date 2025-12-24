package request

type RoleListReq struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
	Desc string `json:"desc" form:"desc"`
	// Keyword  string `json:"keyword" form:"keyword"`
	// PageNum  int    `json:"pageNum" form:"pageNum"`
	// PageSize int    `json:"pageSize" form:"pageSize"`
}

type RoleAddReq struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
	// Type string `json:"type" validate:"required"`
	Desc string `json:"desc"`
}

type RoleDeleteReq struct {
	ID uint `uri:"id" validate:"required"`
}
