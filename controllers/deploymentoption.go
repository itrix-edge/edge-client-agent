package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itrix-edge/edge-client-agent/models"
)

type DeploymentOptionController struct{}

var deploymentOptionModel = new(models.DeploymentOptionModel)

func (m DeploymentOptionController) ListDeploymentOptions(c *gin.Context) {
	listDeployOptions := deploymentOptionModel.ListDeploymentOptions()
	c.JSON(http.StatusOK, gin.H{"data": listDeployOptions})
}

func (m DeploymentOptionController) CreateDeploymentOption(c *gin.Context) {
	var deployOption models.DeploymentOption
	if err := c.ShouldBindJSON(&deployOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	createdDeployOption := deploymentOptionModel.CreateDeploymentOption(&deployOption)
	c.JSON(http.StatusOK, gin.H{"data": createdDeployOption})
}

func (m DeploymentOptionController) ReadDeploymentOptionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	readDeployOption := deploymentOptionModel.GetDeploymentOptionByID(uint(id))
	c.JSON(http.StatusOK, gin.H{"data": readDeployOption})

}

func (m DeploymentOptionController) UpdateDeploymentOptionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var deployOption models.DeploymentOption
	if err := c.ShouldBindJSON(&deployOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	updateDeployOption := deploymentOptionModel.UpdateDeploymentOptionByID(uint(id), &deployOption)
	c.JSON(http.StatusOK, gin.H{"data": updateDeployOption})
}

func (m DeploymentOptionController) DeleteDeploymentOptionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	status := deploymentOptionModel.DeleteDeploymentOptionByID(uint(id))
	c.JSON(http.StatusOK, gin.H{"result": status})
}

func (m DeploymentOptionController) MigrateDeploymentOption(c *gin.Context) {
	deploymentOptionModel.Migrate()
	c.JSON(http.StatusOK, gin.H{"result": true})
}
