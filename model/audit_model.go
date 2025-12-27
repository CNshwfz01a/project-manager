package model

import (
	"project-manager/pkg"
)

type Audit struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Content   string `json:"content" gorm:"type:text;not null"`
	CreatedAt int64  `json:"created_at"`
}

type AuditModel struct{}

var AuditData = &AuditModel{}

// Create 创建审计日志
func (m *AuditModel) Create(content string) error {
	audit := &Audit{
		Content: content,
	}
	return pkg.DB.Create(audit).Error
}

// List 查询审计日志列表
func (m *AuditModel) List(keyword string, startAt, endAt int64, orderBy string, page, pageSize int) ([]Audit, int64, error) {
	var audits []Audit
	var total int64

	query := pkg.DB.Model(&Audit{})

	// 关键字搜索
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 时间范围筛选
	if startAt > 0 {
		query = query.Where("created_at >= ?", startAt)
	}
	if endAt > 0 {
		query = query.Where("created_at <= ?", endAt)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if orderBy == "" {
		orderBy = "created_at"
	}
	query = query.Order(orderBy + " DESC")

	// 分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err := query.Find(&audits).Error
	return audits, total, err
}
