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
		return nil, pkg.NewUnauthorizedError()
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
		isLeader, repError := model.TeamData.IsTeamLeader(r.UserID, r.TeamID)
		if repError != nil {
			return nil, repError
		}
		if !isLeader {
			return nil, pkg.NewUnauthorizedError()
		}
		//判断user是否在team内
		isInTeam, repError := model.TeamData.IsUserInTeam(r.UserID, r.TeamID)
		if repError != nil {
			return nil, repError
		}
		if !isInTeam {
			return nil, pkg.NewUnauthorizedError()
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
	err = model.TeamData.AddUserToTeam(r.UserID, r.TeamID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("添加用户到团队失败: %v", err))
	}
	return nil, nil
}

// AddProjectToTeam ...
func (s *TeamService) AddProjectToTeam(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.TeamAddProjectReq)
	if !ok {
		return nil, ReqAssertErr
	}
	//获取权限 admin可以添加任意project team leader可以添加本team的project 其他403
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		isLeader, repError := model.TeamData.IsTeamLeader(r.TeamID, r.TeamID)
		if repError != nil {
			return nil, repError
		}
		if !isLeader {
			return nil, pkg.NewUnauthorizedError()
		}
		//判断当前team_id是否为leader所属team
		UserID := c.GetUint("user_id")
		isInTeam, repError := model.TeamData.IsUserInTeam(UserID, r.TeamID)
		if repError != nil {
			return nil, repError
		}
		if !isInTeam {
			return nil, pkg.NewUnauthorizedError()
		}
	}
	//查询team是否存在
	_, err := model.TeamData.GetByID(r.TeamID)
	if err == gorm.ErrRecordNotFound {
		return nil, pkg.NewNotFoundError()
	}
	//在TeamModel执行创建project和关联team逻辑
	err = model.TeamData.AddProjectToTeam(r.TeamID, r.ProjectName, r.ProjectDesc)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("添加项目到团队失败: %v", err))
	}
	return nil, nil
}

// Patch ...
func (s *TeamService) Patch(c *gin.Context, teamID uint, req any) (data any, repError any) {
	r, ok := req.(*request.TeamPatch)
	if !ok {
		return nil, ReqAssertErr
	}
	//先获取team
	teamObj, err := model.TeamData.GetByID(teamID)
	if err == gorm.ErrRecordNotFound {
		return nil, pkg.NewNotFoundError()
	} else if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取团队失败: %v", err))
	}
	//根据op path 修改team属性
	for _, opItem := range *r {
		switch opItem.Op {
		case "replace":
			switch opItem.Path {
			case "/leader":
				err = ReplaceLeader(&teamObj, opItem.Value, c)
				if err != nil {
					return nil, err
				}
			case "/name":
				// ReplaceName($teamObj, opItem.Value)
			default:
				//do nothing
			}
		default:
			//do nothing
		}
	}
	//保存teamObj
	err = pkg.DB.Save(&teamObj).Error
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("保存团队失败: %v", err))
	}
	return teamObj, nil
}

// ReplaceLeader ...
func ReplaceLeader(teamObj **model.Team, value any, c *gin.Context) error {
	//判断权限 admin可以修改任意team leader可以修改本team的其他user 其他403
	isAdmin, adminErr := model.GetRoleByName("admin", c)
	if adminErr != nil {
		return adminErr.(error)
	}
	//判断leaderID和登录用户的id是否相同
	UserID := c.GetUint("user_id")
	if !isAdmin {
		isLeader, repError := model.TeamData.IsTeamLeader(UserID, (*teamObj).ID)
		if repError != nil {
			return repError
		}
		if !isLeader {
			return pkg.NewUnauthorizedError()
		}
	}
	//判断value是否为null
	if value == nil {
		(*teamObj).LeaderID = nil
	} else {
		valMap, ok := value.(map[string]interface{})
		if !ok {
			//返回error类型错误
			return pkg.NewRspError(500, fmt.Errorf("无效的负责人值"))
		}
		idFloat, ok := valMap["id"].(float64)
		if !ok {
			return pkg.NewRspError(500, fmt.Errorf("无效的负责人值"))
		}
		leaderID := uint(idFloat)
		if (UserID == leaderID) && !isAdmin {
			//400
			return pkg.NewRspError(400, fmt.Errorf("不可修改自己为团队负责人"))
		}
		//检查leaderID是否存在
		_, err := model.UserData.GetByID(leaderID)
		if err == gorm.ErrRecordNotFound {
			return pkg.NewNotFoundError()
		} else if err != nil {
			return err
		}
		//判断leaderID是否在team内
		isInTeam, repError := model.TeamData.IsUserInTeam(leaderID, (*teamObj).ID)
		if repError != nil {
			return repError
		}
		if !isInTeam {
			//404
			return pkg.NewNotFoundError()
		}
	}

	return nil
}
