package controllers

import (
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
)

var (
	Role = &RoleController{}
)

func Handle(c *gin.Context, req any, fn func() (any, any)) {
	// var err error
	if err := c.Bind(req); err != nil {
		pkg.Err(c, pkg.NewValidatorError(err), nil)
		return
	}

	data, err1 := fn()
	if err1 != nil {
		pkg.Err(c, pkg.ReloadErr(err1), data)
		return
	}
	pkg.Success(c, data)
}
