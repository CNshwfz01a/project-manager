package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditListReq struct {
	OrderBy  string `form:"order_by" validate:"omitempty,oneof=created_at updated_at"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" validate:"omitempty"`
	StartAt  int64  `form:"start_at" validate:"omitempty"`
	EndAt    int64  `form:"end_at" validate:"omitempty"`
}

// SetDefaults 设置默认值
func (r *AuditListReq) SetDefaults() {
	if r.OrderBy == "" {
		r.OrderBy = "created_at"
	}
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
}

// SetParams 设置查询参数
func (r *AuditListReq) SetParams(c *gin.Context) {
	r.OrderBy = c.DefaultQuery("order_by", r.OrderBy)
	r.Page, _ = strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(r.Page)))
	r.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(r.PageSize)))
	r.Keyword = c.DefaultQuery("keyword", r.Keyword)

	startAtStr := c.DefaultQuery("start_at", "")
	if startAtStr != "" {
		r.StartAt, _ = strconv.ParseInt(startAtStr, 10, 64)
	}

	endAtStr := c.DefaultQuery("end_at", "")
	if endAtStr != "" {
		r.EndAt, _ = strconv.ParseInt(endAtStr, 10, 64)
	}
}
