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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct{}

func (s *UserService) Login(c *gin.Context, req any) (data any, repError any) {
	_ = c
	r, ok := req.(*request.UserLoginReq)
	if !ok {
		return nil, ReqAssertErr
	}
	user, err := model.UserData.GetByUsername(r.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NewUnauthorizedError()
		}
		return nil, pkg.NewMySqlError(fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}
	var password = md5.Sum([]byte(r.Password))
	if fmt.Sprintf("%x", password) != user.Password {
		//400
		return nil, pkg.NewRspError(400, fmt.Errorf("密码错误"))
	}
	//cookie写入session
	session := sessions.Default(c)
	sessionId := uuid.New().String()
	session.Set(sessionId, user.ID)
	sec, err := setting.Cfg.GetSection("user")
	if err != nil {
		//返回500
		return nil, pkg.NewRspError(pkg.SystemErr, fmt.Errorf("获取用户配置失败: %s", err.Error()))
	}
	sessionTTL := sec.Key("SESSION_MAX_AGE").MustInt(3600)
	c.SetCookie("session-login", sessionId, sessionTTL, "/", "", false, true)
	session.Options(sessions.Options{
		MaxAge: sessionTTL,
	})
	session.Save()

	return nil, nil
}

// 删除cookie
func (s *UserService) Logout(c *gin.Context) (data any, repError any) {
	session := sessions.Default(c)
	cookie, err := c.Cookie("session-login")
	if err != nil {
		return nil, pkg.NewRspError(400, fmt.Errorf("获取会话失败: %s", err.Error()))
	}
	session.Delete(cookie)
	session.Save()
	c.SetCookie("session-login", "", -1, "/", "", false, true)
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
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
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
