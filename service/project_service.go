package service

import (
	"fmt"
	"project-manager/model"
	"project-manager/model/request"
	"project-manager/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectService struct{}

func (s *ProjectService) AddUserToProject(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.ProjectAddUserReq)
	if !ok {
		return nil, ReqAssertErr
	}

	//验证项目和用户是否存在
	project, err := model.ProjectData.GetByID(r.ProjectID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目信息失败: %v", err))
	}
	if project == nil {
		return nil, pkg.NewNotFoundError()
	}

	targetUser, err := model.UserData.GetByID(r.UserID)
	if err == gorm.ErrRecordNotFound || targetUser == nil {
		return nil, pkg.NewNotFoundError()
	}
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取用户信息失败: %v", err))
	}

	//检查用户是否已在项目中 可以在e2e测试中检查
	isInProject, err := model.ProjectData.IsUserInProject(r.UserID, r.ProjectID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("检查用户项目关系失败: %v", err))
	}
	if isInProject {
		return nil, pkg.NewConflictError(fmt.Errorf("用户已在项目中"))
	}

	//获取项目所属的 Team
	teamID, err := model.ProjectData.GetTeamIDByProjectID(r.ProjectID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目所属团队失败: %v", err))
	}
	if teamID == 0 {
		return nil, pkg.NewRspError(400, fmt.Errorf("项目未关联团队"))
	}

	//权限检查
	currentUserID := c.GetUint("user_id")
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}

	if !isAdmin {
		// 非 admin，检查是否是 Team Leader
		isLeader, err := model.TeamData.IsTeamLeader(currentUserID, teamID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("检查团队领导权限失败: %v", err))
		}
		if !isLeader {
			// 既不是 admin 也不是 Team Leader，无权限
			return nil, pkg.NewUnauthorizedError()
		}

		// 是 Team Leader，检查目标用户是否可见
		// Team Leader 可以看到他 lead 的所有 team 中的用户
		teamsLedByUser, err := model.TeamData.GetTeamsLedByUser(currentUserID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("获取领导的团队失败: %v", err))
		}

		var teamIDs []uint
		for _, team := range teamsLedByUser {
			teamIDs = append(teamIDs, team.ID)
		}

		visibleUserIDs, err := model.TeamData.GetUsersInTeams(teamIDs)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("获取可见用户失败: %v", err))
		}

		// 检查目标用户是否在可见用户列表中
		isVisible := false
		for _, visibleUserID := range visibleUserIDs {
			if visibleUserID == r.UserID {
				isVisible = true
				break
			}
		}

		if !isVisible {
			return nil, pkg.NewUnauthorizedError()
		}
	}

	//添加用户到项目
	err = model.ProjectData.AddUserToProject(r.UserID, r.ProjectID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("添加用户到项目失败: %v", err))
	}

	//如果用户还不是 Team 的成员，自动加入 Team
	isInTeam, err := model.TeamData.IsUserInTeam(r.UserID, teamID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("检查用户团队关系失败: %v", err))
	}

	if !isInTeam {
		err = model.TeamData.AddUserToTeam(r.UserID, teamID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("自动添加用户到团队失败: %v", err))
		}
	}

	return nil, nil
}

/*
*
admin 可以更新任何 Project。
Team Leader 可以更新其 Team 下的 Project。
普通用户无权限更新 Project。
*/
func (s *ProjectService) UpdateProject(c *gin.Context, projectID int, req any) (data any, repError any) {
	r, ok := req.(*request.ProjectUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	//验证项目是否存在
	project, err := model.ProjectData.GetByID(uint(projectID))
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目信息失败: %v", err))
	}
	if project == nil {
		return nil, pkg.NewNotFoundError()
	}
	//权限检查
	currentUserID := c.GetUint("user_id")
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		// 非 admin，检查是否是 Team Leader
		teamID, err := model.ProjectData.GetTeamIDByProjectID(uint(projectID))
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("获取项目所属团队失败: %v", err))
		}
		if teamID == 0 {
			return nil, pkg.NewRspError(400, fmt.Errorf("项目未关联团队"))
		}
		isLeader, err := model.TeamData.IsTeamLeader(currentUserID, teamID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("检查团队领导权限失败: %v", err))
		}
		if !isLeader {
			// 既不是 admin 也不是 Team Leader，无权限
			return nil, pkg.NewUnauthorizedError()
		}
	}
	//更新项目
	project.Name = r.Name
	//判断重名
	exist, err := model.ProjectData.GetByName(r.Name)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目信息失败: %v", err))
	}
	if exist != nil && exist.ID != project.ID {
		return nil, pkg.NewConflictError(fmt.Errorf("资源已存在或存在冲突"))
	}
	project.Status = r.Status
	if r.Desc != nil {
		project.Desc = r.Desc
	}
	err = pkg.DB.Debug().Save(project).Error
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("更新项目信息失败: %v", err))
	}
	return nil, nil
}
