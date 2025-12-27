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

// UpdateProject 更新项目信息
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

// PartialUpdateProject 部分更新项目信息
func (s *ProjectService) PartialUpdateProject(c *gin.Context, projectID int, req any) (data any, repError any) {
	r, ok := req.(*request.ProjectPatch)
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

	//根据op path 修改project属性
	for _, opItem := range *r {
		switch opItem.Op {
		case "replace":
			switch opItem.Path {
			case "/name":
				if name, ok := opItem.Value.(string); ok {
					//判断重名
					exist, err := model.ProjectData.GetByName(name)
					if err != nil {
						return nil, pkg.NewRspError(500, fmt.Errorf("获取项目信息失败: %v", err))
					}
					if exist != nil && exist.ID != project.ID {
						return nil, pkg.NewConflictError(fmt.Errorf("资源已存在或存在冲突"))
					}
					project.Name = name
				} else {
					return nil, pkg.NewRspError(400, fmt.Errorf("无效的项目名称"))
				}
			case "/desc":
				if opItem.Value == nil {
					project.Desc = nil
				} else if desc, ok := opItem.Value.(string); ok {
					project.Desc = &desc
				} else {
					return nil, pkg.NewRspError(400, fmt.Errorf("无效的项目描述"))
				}
			case "/status":
				if status, ok := opItem.Value.(string); ok {
					// 验证状态值
					if status != "WAIT_FOR_SCHEDULE" && status != "IN_PROGRESS" && status != "FINISHED" {
						return nil, pkg.NewRspError(400, fmt.Errorf("无效的项目状态"))
					}
					project.Status = status
				} else {
					return nil, pkg.NewRspError(400, fmt.Errorf("无效的项目状态"))
				}
			default:
				//do nothing
			}
		default:
			//do nothing
		}
	}

	//保存project
	err = pkg.DB.Save(&project).Error
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("保存项目失败: %v", err))
	}
	return project, nil
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(c *gin.Context, projectID int) (data any, repError any) {
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
	//删除项目
	err = model.ProjectData.DeleteProject(uint(projectID))
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("删除项目失败: %v", err))
	}
	return nil, nil
}

// RemoveUserFromProject 清退项目中的用户
func (s *ProjectService) RemoveUserFromProject(c *gin.Context, projectID int, userID int) (data any, repError any) {
	//验证项目是否存在
	project, err := model.ProjectData.GetByID(uint(projectID))
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目信息失败: %v", err))
	}
	if project == nil {
		return nil, pkg.NewNotFoundError()
	}
	//验证用户是否存在
	targetUser, err := model.UserData.GetByID(uint(userID))
	if err == gorm.ErrRecordNotFound || targetUser == nil {
		return nil, pkg.NewNotFoundError()
	}
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取用户信息失败: %v", err))
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
	//清退用户
	err = model.ProjectData.RemoveUserFromProject(uint(userID), uint(projectID))
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("清退用户失败: %v", err))
	}
	return nil, nil
}

// GetProjectUsers 获取项目中的所有用户
func (s *ProjectService) GetProjectUsers(c *gin.Context, projectID int, req any) (data any, repError any) {
	r, ok := req.(*request.UserListReq)
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
	//权限检查 admin可以查看所有用户 team leader可以查看自己下项目的用户 其他用户查看自己参与的项目用户
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
			// 既不是 admin 也不是 Team Leader，检查是否是项目成员
			isInProject, err := model.ProjectData.IsUserInProject(currentUserID, uint(projectID))
			if err != nil {
				return nil, pkg.NewRspError(500, fmt.Errorf("检查用户项目关系失败: %v", err))
			}
			if !isInProject {
				// 既不是 admin 也不是 Team Leader 也不是项目成员，无权限
				return nil, pkg.NewUnauthorizedError()
			}
		}
	}
	//获取项目用户列表
	users, err := model.ProjectData.GetUsersInProject(uint(projectID), r.OrderBy, r.Page, r.PageSize, r.Name)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目用户列表失败: %v", err))
	}
	//格式化返回
	count := len(users)
	return map[string]any{
		"list":  users,
		"count": count,
	}, nil
}

// GetProjectDetail 获取项目详情
func (s *ProjectService) GetProjectDetail(c *gin.Context, projectID int) (data any, repError any) {
	//验证项目是否存在
	project, err := model.ProjectData.GetByID(uint(projectID))
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目信息失败: %v", err))
	}
	if project == nil {
		return nil, pkg.NewNotFoundError()
	}
	//权限检查 admin可以任意project team leader可以查看自己team下的项目 其他用户查看自己参与的项目
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
			// 既不是 admin 也不是 Team Leader，检查是否是项目成员
			isInProject, err := model.ProjectData.IsUserInProject(currentUserID, uint(projectID))
			if err != nil {
				return nil, pkg.NewRspError(500, fmt.Errorf("检查用户项目关系失败: %v", err))
			}
			if !isInProject {
				// 既不是 admin 也不是 Team Leader 也不是项目成员，无权限
				return nil, pkg.NewUnauthorizedError()
			}
		}
	}
	return project, nil
}
