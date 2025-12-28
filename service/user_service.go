package service

import (
	"crypto/md5"
	"fmt"

	// "log"
	"project-manager/model"
	"project-manager/model/request"
	"project-manager/model/response"
	"project-manager/pkg"
	"project-manager/pkg/setting"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct{}

func (s *UserService) Login(c *gin.Context, req any) (data any, repError any) {
	_ = c
	r, ok := req.(*request.UserLoginReq)
	if !ok {
		return nil, ReqAssertErr
	}
	var user *model.User
	var err error
	if r.Username != "" {
		user, err = model.UserData.GetByUsername(r.Username)
	} else if r.Email != "" {
		user, err = model.UserData.GetByEmail(r.Email)
	} else {
		return nil, pkg.NewRspError(400, fmt.Errorf("用户名或邮箱不能为空"))
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewUnauthorizedError()
		}
		return nil, pkg.NewMySqlError(fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	var password = md5.Sum([]byte(r.Password))
	if fmt.Sprintf("%x", password) != user.Password {
		//401
		return nil, pkg.NewNotlogInError()
	}
	//cookie写入session
	session := sessions.Default(c)
	// 设置session过期时间
	sec, err := setting.Cfg.GetSection("user")
	if err != nil {
		//返回500
		return nil, pkg.NewRspError(pkg.SystemErr, fmt.Errorf("获取用户配置失败: %s", err.Error()))
	}
	sessionTTL := sec.Key("SESSION_MAX_AGE").MustInt(3600)

	session.Options(sessions.Options{
		MaxAge: sessionTTL,
		Path:   "/",
	})
	// 直接在Session中存储用户ID
	session.Set("user_id", user.ID)
	session.Save()

	return nil, nil
}

// 删除cookie
func (s *UserService) Logout(c *gin.Context) (data any, repError any) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()
	return nil, nil
}

func (s *UserService) ChangePassword(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.UserChangePasswordReq)
	if !ok {
		return nil, ReqAssertErr
	}
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)
	user, err := model.UserData.GetByID(userID)
	if err != nil {
		return nil, pkg.NewRspError(400, fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	var currentPasswordHash = md5.Sum([]byte(r.OldPassword))
	if fmt.Sprintf("%x", currentPasswordHash) != user.Password {
		return nil, pkg.NewRspError(400, fmt.Errorf("旧密码不正确"))
	}
	var newPasswordHash = md5.Sum([]byte(r.NewPassword))
	user.Password = fmt.Sprintf("%x", newPasswordHash)
	//新密码不能与旧密码相同
	if user.Password == fmt.Sprintf("%x", currentPasswordHash) {
		return nil, pkg.NewRspError(400, fmt.Errorf("新密码不能与旧密码相同"))
	}
	//修改用户状态为激活
	user.Status = 1
	err = pkg.DB.Save(user).Error
	if err != nil {
		return nil, pkg.NewRspError(400, fmt.Errorf("更新用户密码失败: %s", err.Error()))
	}
	//调用登出方法
	_, repError = s.Logout(c)
	if repError != nil {
		return nil, repError
	}
	return nil, nil
}

/*
*
只有admin用户可以执行操作
创建时指定用户名和初始密码
自动绑定 normal user Role
*/
func (s *UserService) Add(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.UserAddReq)
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
	_, err := model.UserData.GetByUsername(r.Username)
	if err == nil {
		return nil, pkg.NewConflictError(fmt.Errorf("资源已存在或存在冲突"))
	}
	//创建用户
	var passwordHash = md5.Sum([]byte(r.Password))
	user := &model.User{
		Username: r.Username,
		Password: fmt.Sprintf("%x", passwordHash),
	}
	//设置nickname为username
	user.Nickname = &r.Username
	err = pkg.DB.Create(user).Error
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("创建用户失败: %s", err.Error()))
	}
	//绑定normal user 角色
	role, err := model.RoleData.GetByName("normal user")
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取角色信息失败: %s", err.Error()))
	}
	err = pkg.DB.Model(user).Association("Roles").Append(role)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("绑定角色失败: %s", err.Error()))
	}
	//返回user详情
	return formatUserDetail(user), nil
}

// format user detail
func formatUserDetail(user *model.User) *response.UserListResp {
	userDetail := response.UserListResp{
		ID:        user.ID,
		Username:  user.Username,
		Status:    user.Status,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Logo:      user.Logo,
		Roles:     []request.RoleListReq{},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	var roleList []request.RoleListReq
	for _, role := range user.Roles {
		roleList = append(roleList, request.RoleListReq{
			ID:   role.ID,
			Name: role.Name,
			Type: role.Type,
			Desc: role.Desc,
		})
	}
	userDetail.Roles = roleList
	return &userDetail
}

func (s *UserService) getUserDetail(userID uint) (userDetail *response.UserListResp, repError any) {
	user, err := model.UserData.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//404
			return nil, pkg.NewNotFoundError()
		}
		return nil, pkg.NewRspError(400, fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	return formatUserDetail(user), nil
}

func (s *UserService) List(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.UserListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	//调用查询接口获取结果
	users, err := model.UserData.List(c, r)
	if err != nil {
		return nil, pkg.NewRspError(400, fmt.Errorf("获取用户列表失败: %s", err.Error()))
	}
	//formatUserList
	respList := formatUserList(users)
	//计算respList的容量
	var count int64 = int64(len(respList))
	//返回{list total}
	return map[string]any{
		"list":  respList,
		"total": count,
	}, nil

}

// 将格式化用户信息的功能封装成函数
func formatUserList(users []model.User) []response.UserListResp {
	var respList []response.UserListResp
	for _, user := range users {
		formatUserDetail(&user)
		respList = append(respList, *formatUserDetail(&user))
	}
	return respList
}

// Delete 删除用户
func (s *UserService) Delete(c *gin.Context, userID uint) (data any, repError any) {
	//获取权限
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		//403
		return nil, pkg.NewUnauthorizedError()
	}
	//admin不能被删除 有风险 但是先这么写
	if userID == 1 {
		return nil, pkg.NewRspError(400, fmt.Errorf("管理员用户不能被删除"))
	}
	//删除用户
	err := model.UserData.DeleteByID(userID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("删除用户失败: %s", err.Error()))
	}
	return nil, nil
}

// 分配用户角色
func (s *UserService) AssignRole(c *gin.Context, req any, userID uint) (data any, repError any) {
	r, ok := req.(*request.UserAssignRoleReq)
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
	//获取用户
	user, err := model.UserData.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewRspError(400, fmt.Errorf("用户不存在"))
		}
		return nil, pkg.NewRspError(400, fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	//获取角色
	role, err := model.RoleData.GetByID(r.RoleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewRspError(400, fmt.Errorf("角色不存在"))
		}
		return nil, pkg.NewRspError(400, fmt.Errorf("获取角色信息失败: %s", err.Error()))
	}
	//不能删除role type为System的角色
	if role.Type == "System" {
		return nil, pkg.NewRspError(400, fmt.Errorf("不能分配系统内置角色"))
	}
	//分配角色
	err = pkg.DB.Model(user).Association("Roles").Append(role)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("分配角色失败: %s", err.Error()))
	}
	return nil, nil
}

// 移除用户角色
func (s *UserService) RemoveRole(c *gin.Context, userID uint, roleID uint) (data any, repError any) {
	//获取权限
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}
	if !isAdmin {
		//403
		return nil, pkg.NewUnauthorizedError()
	}
	//获取用户
	user, err := model.UserData.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewRspError(400, fmt.Errorf("用户不存在"))
		}
		return nil, pkg.NewRspError(400, fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	//获取角色
	role, err := model.RoleData.GetByID(roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewRspError(400, fmt.Errorf("角色不存在"))
		}
		return nil, pkg.NewRspError(400, fmt.Errorf("获取角色信息失败: %s", err.Error()))
	}
	//不能删除role type为system的角色
	if role.Type == "System" {
		return nil, pkg.NewUnauthorizedError()
	}
	//移除角色
	err = pkg.DB.Model(user).Association("Roles").Delete(role)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("移除角色失败: %s", err.Error()))
	}
	return nil, nil
}

// 自身信息
func (s *UserService) MyDetail(c *gin.Context) (data any, repError any) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)
	return s.getUserDetail(userID)
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(c *gin.Context, req any) (data any, repError any) {
	r, ok := req.(*request.UserUpdateProfileReq)
	if !ok {
		return nil, ReqAssertErr
	}
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)
	user, err := model.UserData.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewRspError(400, fmt.Errorf("用户不存在"))
		}
		return nil, pkg.NewRspError(400, fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	user.Email = &r.Email
	user.Nickname = &r.Nickname
	user.Logo = &r.Logo
	err = pkg.DB.Save(user).Error
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("更新用户资料失败: %s", err.Error()))
	}
	return s.getUserDetail(userID)
}

// Detail 查询用户详情
func (s *UserService) Detail(c *gin.Context, userID uint) (data any, repError any) {
	//判断用户是否跟登录人在同team中
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	currentUserID := userIDInterface.(uint)

	// 如果查询的是自己，直接返回
	if currentUserID == userID {
		return s.getUserDetail(userID)
	}

	//获取权限
	isAdmin, repError := model.GetRoleByName("admin", c)
	if repError != nil {
		return nil, repError
	}

	if !isAdmin {
		// 非admin角色 只能查看相同team下的用户
		isSameTeam, err := model.TeamData.IsSameTeam(currentUserID, userID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("查询团队信息失败: %s", err.Error()))
		}
		if !isSameTeam {
			return nil, pkg.NewUnauthorizedError()
		}
	}

	return s.getUserDetail(userID)
}

/**
用户团队列表
查询他人的 teams 列表时首先要满足可见性约束。
admin 用户可以查看所有 users 的所有 teams。
普通用户仅能查看共同 Teams 下用户的 teams 列表。
普通用户能查看其他用户的 teams 列表不应当超过 Me 自身的列表范围,比如 Me 所在的 teams 为 [a, b], user 所在的 teams 为 [a, c], 那么 Me 查询该接口仅返回 [a]。
*/

func (s *UserService) TeamList(c *gin.Context, userID uint, req *request.UserListReq) (data any, repError any) {
	// Get current user
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	currentUserID := userIDInterface.(uint)

	// Check if admin
	isAdmin, err := model.GetRoleByName("admin", c)
	if err != nil {
		return nil, err
	}

	var teams []model.Team
	var dbErr error

	if isAdmin {
		teams, dbErr = model.TeamData.GetTeamsByUserID(userID, req)
		if dbErr != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("获取团队列表失败: %s", dbErr.Error()))
		}
	} else {
		// Normal user: get common teams
		teams, dbErr = model.TeamData.GetCommonTeams(currentUserID, userID, req)
	}

	if dbErr != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取团队列表失败: %s", dbErr.Error()))
	}

	total := len(teams)

	// If leading is true, total should reflect the number of teams led by the target user

	// Format response
	return map[string]any{
		"list":  formatTeamList(teams),
		"total": total,
	}, nil
}

func formatTeamList(teams []model.Team) []response.TeamGetResp {
	var respList []response.TeamGetResp
	for _, team := range teams {
		var desc string
		if team.Desc != nil {
			desc = *team.Desc
		}
		respList = append(respList, response.TeamGetResp{
			ID:        team.ID,
			Name:      team.Name,
			Desc:      desc,
			Leader:    team.Leader,
			Projects:  team.Projects,
			CreatedAt: team.CreatedAt,
			UpdatedAt: team.UpdatedAt,
		})
	}
	// return empty slice instead of nil if empty
	if respList == nil {
		return []response.TeamGetResp{}
	}
	return respList
}

/*
*
用户项目列表
admin 可以查看所有用户的所有参与项目的列表。
普通用户仅可以查看同 Team 的用户参与的项目列表。
普通用户能查看到的其他用户的 projects 列表范围不应当超出自身参与的 projects 范围。
*/
func (s *UserService) ProjectList(c *gin.Context, userID uint) (data any, repError any) {
	// Get current user
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	currentUserID := userIDInterface.(uint)

	// Check if admin
	isAdmin, err := model.GetRoleByName("admin", c)
	if err != nil {
		return nil, err
	}

	var projects []model.Project
	var dbErr error

	if isAdmin || currentUserID == userID {
		// Admin or self: get all projects of the target user
		projects, dbErr = model.ProjectData.GetProjectsByUserID(userID)
	} else {
		// Normal user: check same team first
		isSameTeam, err := model.TeamData.IsSameTeam(currentUserID, userID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("查询团队信息失败: %s", err.Error()))
		}
		if !isSameTeam {
			return nil, pkg.NewUnauthorizedError()
		}

		// Get common projects
		projects, dbErr = model.ProjectData.GetCommonProjects(currentUserID, userID)
	}

	if dbErr != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取项目列表失败: %s", dbErr.Error()))
	}

	// Format response
	return map[string]any{
		"list":  formatProjectList(projects),
		"total": len(projects),
	}, nil
}

func formatProjectList(projects []model.Project) []response.ProjectListResp {
	var respList []response.ProjectListResp
	for _, project := range projects {
		var desc string
		if project.Desc != nil {
			desc = *project.Desc
		}
		respList = append(respList, response.ProjectListResp{
			ID:        project.ID,
			Name:      project.Name,
			Desc:      desc,
			Status:    project.Status,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
		})
	}
	if respList == nil {
		return []response.ProjectListResp{}
	}
	return respList
}

// MyTeamList 我所在的团队列表
func (s *UserService) MyTeamList(c *gin.Context, leading int) (data any, repError any) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)
	teams, err := model.TeamData.GetByTeamList(userID, leading)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("获取我的团队列表失败: %s", err.Error()))
	}
	return map[string]any{
		"list":  formatTeamList(teams),
		"total": len(teams),
	}, nil

}

// MyProjectList 我所在的项目列表
func (s *UserService) MyProjectList(c *gin.Context) (data any, repError any) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)
	return s.ProjectList(c, userID)
}

// LeaveTeam 退出团队
func (s *UserService) LeaveTeam(c *gin.Context, teamID uint) (data any, repError any) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)

	// 检查是否在团队中
	inTeam, err := model.TeamData.IsUserInTeam(userID, teamID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("查询团队成员失败: %s", err.Error()))
	}
	if !inTeam {
		return nil, pkg.NewRspError(400, fmt.Errorf("你不在该团队中"))
	}

	// 检查是否是Leader 是leader则将team的leader_id设为null
	isLeader, err := model.TeamData.IsTeamLeader(userID, teamID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("查询团队信息失败: %s", err.Error()))
	}
	if isLeader {
		err = model.TeamData.RemoveTeamLeader(teamID)
		if err != nil {
			return nil, pkg.NewRspError(500, fmt.Errorf("移除团队负责人失败: %s", err.Error()))
		}
	}

	err = model.TeamData.RemoveUserFromTeam(teamID, userID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("退出团队失败: %s", err.Error()))
	}
	return nil, nil
}

// LeaveProject 退出项目
func (s *UserService) LeaveProject(c *gin.Context, projectID uint) (data any, repError any) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil, pkg.NewRspError(400, fmt.Errorf("未获取到用户登录信息"))
	}
	userID := userIDInterface.(uint)

	// 检查是否在项目中
	inProject, err := model.ProjectData.IsUserInProject(userID, projectID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("查询项目成员失败: %s", err.Error()))
	}
	if !inProject {
		return nil, pkg.NewRspError(400, fmt.Errorf("你不在该项目中"))
	}

	err = model.ProjectData.RemoveUserFromProject(userID, projectID)
	if err != nil {
		return nil, pkg.NewRspError(500, fmt.Errorf("退出项目失败: %s", err.Error()))
	}
	return nil, nil
}
