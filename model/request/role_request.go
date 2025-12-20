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
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	Desc string `json:"desc"`
}
