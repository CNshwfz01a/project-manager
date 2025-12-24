package service

import (
	"fmt"
	"project-manager/model"
	"project-manager/model/request"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeamService struct{}

func (s *TeamService) Add(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.TeamAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	//获取权限
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		//403
		return nil, pkg.NewUnauthorizedError(fmt.Errorf("没有权限执行该操作"))
	}
	//判断重名
	exist, err := model.TeamData.GetByTeamName(r.Name)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, pkg.NewConflictError(fmt.Errorf("资源已存在或存在冲突"))
	}
	//创建team
	newTeam := &model.Team{
		Name: r.Name,
		Desc: r.Desc,
	}
	err = pkg.DB.Create(newTeam).Error
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("创建团队失败: %v", err))
	}
	//增加一个空的leader结构体
	newTeam.Leader = &model.User{}
	return newTeam, nil
}

// AddUserToTeam ...
func (s *TeamService) AddUserToTeam(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.TeamAddUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	//获取权限 admin可以添加任意user team leader可以添加本team的user 其他403
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		isLeader, repError := model.TeamData.IsTeamLeader(uint(r.UserID), r.TeamID)
		if repError != nil {
			return nil, repError
		}
		if !isLeader {
			return nil, pkg.NewUnauthorizedError(fmt.Errorf("没有权限执行该操作"))
		}
		//判断user是否在team内
		isInTeam, repError := model.TeamData.IsUserInTeam(r.UserID, r.TeamID)
		if repError != nil {
			return nil, repError
		}
		if !isInTeam {
			return nil, pkg.NewUnauthorizedError(fmt.Errorf("没有权限执行该操作"))
		}
	}
	//查询team和user是否存在
	_, err := model.TeamData.GetByID(r.TeamID)
	if err == gorm.ErrRecordNotFound {
		return nil, pkg.NewNotFoundError()
	}
	_, err = model.UserData.GetByID(r.UserID)
	if err == gorm.ErrRecordNotFound {
		return nil, pkg.NewNotFoundError()
	}
	//添加user到team
	err := model.TeamData.AddUserToTeam(r.UserID, r.TeamID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("添加用户到团队失败: %v", err))
	}
	return nil, nil
}
