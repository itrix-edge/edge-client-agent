package controllers

import (
	"models"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 *  NamespaceController provides k8s namespace controller
 */
type NamespaceController struct{}

var namespaceModel = new(models.NamespaceModel)

/**
 *  GetNamespaces get all namespaces in the cluster
 */
func (nc NamespaceController) GetNamespaces(c *gin.Context) {

	list, err := namespaceModel.GetNamespaces()
	// user, token, err := userModel.Login(loginForm)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"data": list})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Cannot get namespaces", "error": err.Error()})
	}
}
