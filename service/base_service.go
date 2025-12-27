package service

import (
	"fmt"
	"project-manager/pkg"
)

var (
	ReqAssertErr = pkg.NewRspError(pkg.SystemErr, fmt.Errorf("请求异常"))
	Role         = &RoleService{}
	User         = &UserService{}
	Team         = &TeamService{}
	Project      = &ProjectService{}
	Audit        = &AuditService{}
)
