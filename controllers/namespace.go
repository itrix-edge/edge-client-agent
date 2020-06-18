package controllers

import (
	"net/http"

	"github.com/stevennick/edge-client-agent/models"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	options := v1.ListOptions{}
	list, err := namespaceModel.GetNamespaces(options)
	// user, token, err := userModel.Login(loginForm)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"data": list})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Cannot get namespaces", "error": err.Error()})
	}
}
