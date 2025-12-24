package controllers

import (
	"project-manager/model/request"
	"project-manager/service"

	// "strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (m *UserController) Login(c *gin.Context) {
	req := new(request.UserLoginReq)
	Handle(c, req, func() (any, any) {
		return service.User.Login(c, req)
	})
}

func (m *UserController) Logout(c *gin.Context) {
	Handle(c, nil, func() (any, any) {
		return service.User.Logout(c)
	})
}

// 修改自身密码 修改完成需要重新登录
func (m *UserController) ChangePassword(c *gin.Context) {
	req := new(request.UserChangePasswordReq)
	Handle(c, req, func() (any, any) {
		return service.User.ChangePassword(c, req)
	})
}

func (m *UserController) Add(c *gin.Context) {
	req := new(request.UserAddReq)
	Handle(c, req, func() (any, any) {
		return service.User.Add(c, req)
	})
}

// 查询用户列表
func (m *UserController) List(c *gin.Context) {
	req := new(request.UserListReq)
	req.SetParams(c)
	Handle(c, req, func() (any, any) {
		return service.User.List(c, req)
	})
}

//用户
