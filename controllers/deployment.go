package controllers

import (
	"net/http"

	"github.com/itrix-edge/edge-client-agent/models"
	// "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/deployment"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
 *  DeploymentController provides k8s deployment controller
 */ //
type DeploymentController struct{}

var deploymentModel = new(models.DeploymentModel)

/**
 *  GetDeployments get all Deployments in the cluster
 */
func (nc DeploymentController) GetDeployments(c *gin.Context) {
	namespace := c.Param("namespace")
	listOptions := v1.ListOptions{}
	list, err := deploymentModel.GetDeployments(namespace, listOptions)
	// user, token, err := userModel.Login(loginForm)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"data": list})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Cannot get deployment", "error": err.Error()})
	}
}

// CreateNewDeployment create new deployment. Currently this only create under default namespace
func (nc DeploymentController) CreateNewDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	deploymentResult, err := deploymentModel.CreateDeployment(namespace, &deployment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": deploymentResult})
}

func (nc DeploymentController) ReadDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var getOptions = v1.GetOptions{}
	deploymentResult, err := deploymentModel.ReadDeployment(namespace, name, getOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": deploymentResult})
}

func (nc DeploymentController) UpdateDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	deploymentResult, err := deploymentModel.UpdateDeplyment(namespace, &deployment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": deploymentResult})
}

func (nc DeploymentController) DeleteDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	deploymentResult := deploymentModel.DeleteDeployment(namespace, deployment.Name, &v1.DeleteOptions{})
	c.JSON(http.StatusOK, gin.H{"result": deploymentResult})
}

func (nc DeploymentController) DeleteDeploymentByName(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// var deployment appsv1.Deployment
	// if err := c.ShouldBindJSON(&deployment); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// }
	deployment, err := deploymentModel.ReadDeployment(namespace, name, v1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	status := deploymentModel.DeleteDeployment(namespace, deployment.Name, &v1.DeleteOptions{})
	c.JSON(http.StatusOK, gin.H{"result": status})
}
