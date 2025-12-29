package controllers

import (
	"project-manager/model/request"
	"project-manager/service"

	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamController struct{}

func (m *TeamController) Add(c *gin.Context) {
	req := new(request.TeamAddReq)
	Handle(c, req, func() (any, any) {
		return service.Team.Add(c, req)
	})
}

func (m *TeamController) AddUserToTeam(c *gin.Context) {
	req := new(request.TeamAddUserReq)
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	req.TeamID = uint(id)
	Handle(c, req, func() (any, any) {
		return service.Team.AddUserToTeam(c, req)
	})
}

func (m *TeamController) AddProjectToTeam(c *gin.Context) {
	req := new(request.TeamAddProjectReq)
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	req.TeamID = uint(id)
	Handle(c, req, func() (any, any) {
		return service.Team.AddProjectToTeam(c, req)
	})
}

/*
	/type TeamPatch []struct {
		Op    string      `json:"op" validate:"required,oneof=replace add remove"`
		Path  string      `json:"path" validate:"required,oneof=/name /desc /leader"`
		Value interface{} `json:"value,omitempty"`
	}

*
*/
func (m *TeamController) Patch(c *gin.Context) {
	req := new(request.TeamPatch)
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	Handle(c, req, func() (any, any) {
		return service.Team.Patch(c, uint(id), req)
	})
}

func (m *TeamController) Get(c *gin.Context) {
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	Handle(c, nil, func() (any, any) {
		return service.Team.Get(c, uint(id))
	})
}

func (m *TeamController) List(c *gin.Context) {
	req := new(request.UserListReq) //复用UserListReq的分页排序参数
	Handle(c, req, func() (any, any) {
		return service.Team.List(c, req)
	})
}

func (m *TeamController) ListUsers(c *gin.Context) {
	req := new(request.UserListReq) //复用UserListReq的分页排序参数
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	Handle(c, req, func() (any, any) {
		return service.Team.ListUsers(c, uint(id), req)
	})
}

// Update ...
func (m *TeamController) Update(c *gin.Context) {
	req := new(request.TeamUpdateReq)
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	Handle(c, req, func() (any, any) {
		return service.Team.Update(c, uint(id), req)
	})
}

// Delete ...
func (m *TeamController) Delete(c *gin.Context) {
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	Handle(c, nil, func() (any, any) {
		return service.Team.Delete(c, uint(id))
	})
}

// RemoveUserFromTeam ...
func (m *TeamController) RemoveUserFromTeam(c *gin.Context) {
	//读url的team_id参数
	var teamID, _ = strconv.Atoi(c.Param("team_id"))
	//读url的user_id参数
	var userID, _ = strconv.Atoi(c.Param("user_id"))
	Handle(c, nil, func() (any, any) {
		return service.Team.RemoveUserFromTeam(c, uint(userID), uint(teamID))
	})
}

// ListProjects ...
func (m *TeamController) ListProjects(c *gin.Context) {
	req := new(request.ProjectListReq) //复用ProjectListReq的分页排序参数
	//读url的team_id参数
	var id, _ = strconv.Atoi(c.Param("team_id"))
	req.TeamID = uint(id)
	Handle(c, req, func() (any, any) {
		return service.Team.ListProjects(c, uint(id), req)
	})
}
