package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectAddUserReq struct {
	ProjectID uint `json:"project_id" validate:"required"`
	UserID    uint `json:"user_id" validate:"required"`
}

type ProjectListReq struct {
	OrderBy  string `form:"order_by" validate:"omitempty,oneof=id name created_at updated_at"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
	Name     string `form:"name" validate:"omitempty"`
	TeamID   uint   `form:"team_id" validate:"omitempty"`
	PartIn   bool   `form:"part_in" validate:"omitempty"`
}

// ProjectListReq 默认值
func (r *ProjectListReq) SetDefaults() {
	if r.OrderBy == "" {
		r.OrderBy = "id"
	}
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}
}

// ProjectListReq 设置查询参数
func (r *ProjectListReq) SetParams(c *gin.Context) {
	r.OrderBy = c.DefaultQuery("order_by", r.OrderBy)
	r.Page, _ = strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(r.Page)))
	r.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(r.PageSize)))
	r.Name = c.DefaultQuery("name", r.Name)
	teamIDStr := c.DefaultQuery("team_id", "")
	if teamIDStr != "" {
		teamID, _ := strconv.Atoi(teamIDStr)
		r.TeamID = uint(teamID)
	}
	partInStr := c.DefaultQuery("part_in", "")
	if partInStr != "" {
		partIn, _ := strconv.ParseBool(partInStr)
		r.PartIn = partIn
	}
}

type ProjectUpdateReq struct {
	Name   string  `json:"name" validate:"required"`
	Desc   *string `json:"desc" validate:"omitempty"`
	Status string  `json:"status" validate:"omitempty,oneof=WAIT_FOR_SCHEDULE IN_PROGRESS FINISHED"`
}

func (r *ProjectUpdateReq) SetDefaults() {
	if r.Status == "" {
		r.Status = "WAIT_FOR_SCHEDULE"
	}
}
