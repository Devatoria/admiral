package api

import (
	"net/http"
	"strconv"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
)

// getTeamUsers returns the list of users the given team contains
func getTeamUsers(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Get team
	var team models.Team
	db.Instance().Preload("Users").Where("id = ?", id).Find(&team)
	if team.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	// Sanitize passwords
	for i := range team.Users {
		team.Users[i].Password = "[REDACTED]"
	}

	c.JSON(http.StatusOK, team.Users)
}

// postTeamUsers adds a list of users to the given team
func postTeamUsers(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var team models.Team
	db.Instance().Where("id = ?", id).Find(&team)
	if team.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	var data TeamUsers
	err = c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var users []models.User
	db.Instance().Where("id IN (?)", data.Users).Find(&users)

	team.Users = append(team.Users, users...)
	db.Instance().Save(&team)

	c.Status(http.StatusOK)
}

// deleteTeamUser removes the given user from the given team
func deleteTeamUser(c *gin.Context) {
	teamIDParam := c.Param("teamID")
	teamID, err := strconv.Atoi(teamIDParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Get team and user
	var team models.Team
	db.Instance().Where("id = ?", teamID).Find(&team)
	if team.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	var user models.User
	db.Instance().Where("id = ?", userID).Find(&user)
	if user.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	// Remove association
	db.Instance().Model(&team).Association("Users").Delete(&user)
	c.Status(http.StatusOK)
}
