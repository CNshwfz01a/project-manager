package service

import (
	"fmt"
	"project-manager/model"
	"project-manager/model/request"
	"project-manager/model/response"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		return nil, pkg.NewMySqlError(fmt.Errorf("获取角色列表失败: %s", err.Error()))
	}

	count, err := model.RoleData.Count()
	if err != nil {
		return nil, pkg.NewMySqlError(fmt.Errorf("获取角色总数失败: %s", err.Error()))
	}

	return response.RoleListRsp{
		Roles: roles,
		Total: count,
	}, nil
}

func (s *RoleService) Add(c *gin.Context, req any) (data any, repError any) {
	_ = c
	r, ok := req.(*request.RoleAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if model.RoleData.Exist(pkg.H{"name": r.Name}) {
		return nil, pkg.NewConflictError(fmt.Errorf("角色名已存在"))
	}

	//判断当前角色是否有权限
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		return nil, pkg.NewUnauthorizedError()
	}

	role := model.Role{
		Name: r.Name,
		Type: "Custom",
		Desc: r.Desc,
	}

	err := model.RoleData.Add(&role)
	if err != nil {
		return nil, pkg.NewMySqlError(fmt.Errorf("添加角色失败: %s", err.Error()))
	}

	return role, nil
}

func (s *RoleService) Delete(c *gin.Context, req any) (data any, repError any) {
	_ = c
	r, ok := req.(*request.RoleDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	//判断当前角色是否有权限
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		return nil, pkg.NewUnauthorizedError()
	}
	//不可以删除系统角色
	roleData, err := model.RoleData.GetByID(r.ID)
	role := roleData
	if role.Type == "System" {
		return nil, pkg.NewRspError(400, fmt.Errorf("系统角色不能被删除"))
	}
	if err != nil {
		return nil, pkg.NewMySqlError(fmt.Errorf("查询角色失败: %s", err.Error()))
	}
	//删除对应的用户关联关系
	err = model.UserData.DeleteRelationByRoleID(r.ID)
	if err != nil {
		return nil, pkg.NewMySqlError(fmt.Errorf("删除角色关联用户失败: %s", err.Error()))
	}
	//删除角色
	err = model.RoleData.Delete(r.ID)
	if err != nil {
		//判断是否找到记录
		if err == gorm.ErrRecordNotFound {
			//返回404错误
			return nil, pkg.NewNotFoundError()
		}
		return nil, pkg.NewMySqlError(fmt.Errorf("删除角色失败: %s", err.Error()))
	}
	return nil, nil
}
