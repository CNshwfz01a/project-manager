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
