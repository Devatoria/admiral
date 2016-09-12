package api

import (
	"net/http"
	"strconv"

	"github.com/Devatoria/admiral/auth"
	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/filters"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
)

// Team represents a team form
type Team struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// TeamUsers represents a list of users of a team
type TeamUsers struct {
	Users []int `form:"users" json:"users" binding:"required"`
}

// getTeams returns all teams ordered by name
func getTeams(c *gin.Context) {
	nParam := c.DefaultQuery("n", "25")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var teams []models.Team
	db.Instance().Order("name").Limit(n).Find(&teams)

	c.JSON(http.StatusOK, teams)
}

// getTeam returns the team associated to the given ID
func getTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var team models.Team
	db.Instance().Preload("Owner").Where("id = ?", id).Find(&team)
	if team.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, team)
}

// putTeam creates a new team, ensuring that its name is correct and not used
func putTeam(c *gin.Context) {
	var data Team
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err = filters.ValidateTeam(data.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Exists(db.Instance(), "name", data.Name, &models.Team{}) {
		c.Status(http.StatusConflict)
		return
	}

	owner, err := auth.GetCurrentUser(c)
	if err != nil {
		panic(err)
	}

	team := models.Team{Name: data.Name, Owner: owner}
	db.Instance().Create(&team)

	c.JSON(http.StatusOK, team)
}

// deleteTeam deletes the given team
func deleteTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	db.Instance().Where("id = ?", id).Delete(&models.Team{})
	c.Status(http.StatusOK)
}

// patchTeam renames the given team ensuring the new name is not used and correct
func patchTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Data binding
	var data Team
	err = c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Check team exists
	var team models.Team
	db.Instance().Where("id = ?", id).Find(&team)
	if team.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	// Check new name doesn't exist
	if team.Name != data.Name {
		if db.Exists(db.Instance(), "name", data.Name, &models.Team{}) {
			c.Status(http.StatusConflict)
			return
		}
	}

	// Validate name
	if err = filters.ValidateTeam(data.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team.Name = data.Name
	db.Instance().Save(&team)
	c.JSON(http.StatusOK, team)
}
