package service

import (
	"fmt"
	"project-manager/model"
	"project-manager/model/request"
	"project-manager/model/response"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
)

type RoleService struct{}

func (s *RoleService) List(c *gin.Context, req any) (data any, repError any) {
	_ = c
	r, ok := req.(*request.RoleListReq)
	if !ok {
		return nil, ReqAssertErr
	}

	roles, err := model.RoleData.List(r)
	if err != nil {
		return nil, pkg.NewMySqlError(fmt.Errorf("获取菜单列表失败: %s", err.Error()))
	}

	count, err := model.RoleData.Count()
	if err != nil {
		return nil, pkg.NewMySqlError(fmt.Errorf("获取菜单总数失败: %s", err.Error()))
	}

	return response.RoleListRsp{
		Total: count,
		Roles: roles,
	}, nil
}
