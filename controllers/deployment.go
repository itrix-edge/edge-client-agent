package controllers

import (
	"net/http"

	"github.com/stevennick/edge-client-agent/models"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
 *  DeploymentController provides k8s deployment controller
 */
type DeploymentController struct{}

var deploymentModel = new(models.DeploymentModel)

/**
 *  GetDeployments get all Deployments in the cluster
 */
func (nc DeploymentController) GetDeployments(c *gin.Context) {

	options := v1.GetOptions{}
	list, err := deploymentModel.GetDeployments(options)
	// user, token, err := userModel.Login(loginForm)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"data": list})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Cannot get deployment", "error": err.Error()})
	}
}

// CreateNewDeployment create new deployment. Currently this only create under default namespace
func (nc DeploymentController) CreateNewDeployment(c *gin.Context) {
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	deploymentResult, err := deploymentModel.AddNewDeployment("default", &deployment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": deploymentResult})
}
