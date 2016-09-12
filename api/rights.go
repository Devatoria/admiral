package api

import (
	"net/http"
	"strconv"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
)

type TeamNamespaceRight struct {
	TeamID      uint `form:"teamID" json:"teamID" binding:"required"`
	NamespaceID uint `form:"namespaceID" json:"namespaceID" binding:"required"`
	Pull        bool `form:"pull" json:"pull" binding:"required"`
	Push        bool `form:"push" json:"push" binding:"required"`
}

// getTeamNamespaceRights returns the rights for the given team/namespace couple (push, pull, both or nothing)
func getTeamNamespaceRights(c *gin.Context) {
	teamIDParam := c.Param("teamID")
	teamID, err := strconv.Atoi(teamIDParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	namespaceIDParam := c.Param("namespaceID")
	namespaceID, err := strconv.Atoi(namespaceIDParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Get team and namespace rights
	var rights []models.TeamNamespaceRight
	db.Instance().Where("team_id = ? AND namespace_id = ?", teamID, namespaceID).Find(&rights)
	if len(rights) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, rights)
}

// putTeamNamespaceRights grants a right (pull, push or both) to the given team on the given namespace
func putTeamNamespaceRight(c *gin.Context) {
	var data TeamNamespaceRight
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Retrieve team
	var team models.Team
	db.Instance().Where("id = ?", data.TeamID).Find(&team)
	if team.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	// Retrieve namespace
	var namespace models.Namespace
	db.Instance().Where("id = ?", data.NamespaceID).Find(&namespace)
	if namespace.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	// Create right
	right := models.TeamNamespaceRight{
		Team:      team,
		Namespace: namespace,
		Pull:      data.Pull,
		Push:      data.Push,
	}
	db.Instance().Create(&right)

	c.Status(http.StatusOK)
}

// deleteTeamNamespaceRights revokes all rights of the given team to the given namespace
func deleteTeamNamespaceRight(c *gin.Context) {
	teamIDParam := c.Param("teamID")
	teamID, err := strconv.Atoi(teamIDParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	namespaceIDParam := c.Param("namespaceID")
	namespaceID, err := strconv.Atoi(namespaceIDParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var right models.TeamNamespaceRight
	db.Instance().Where("team_id = ? AND namespace_id = ?", teamID, namespaceID).Find(&right)
	if right.TeamID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	db.Instance().Delete(&right)

	c.Status(http.StatusOK)
}
