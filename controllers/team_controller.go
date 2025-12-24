package controllers

import (
	"project-manager/model/request"
	"project-manager/service"

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
	//读uri的team_id参数
	c.BindUri(req)
	Handle(c, req, func() (any, any) {
		return service.Team.AddUserToTeam(c, req)
	})
}
