package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/itrix-edge/edge-client-agent/models"

	"github.com/gin-gonic/gin"
)

type HookController struct{}

var hookModel = new(models.HookModel)

func (m HookController) ListHooks(c *gin.Context) {
	hooks := hookModel.ListHooks()
	c.JSON(http.StatusOK, gin.H{"data": hooks})
}

func (m HookController) CreateHook(c *gin.Context) {
	var hook models.Hook
	if err := c.ShouldBindJSON(&hook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	createdHook := hookModel.CreateHook(&hook)
	c.JSON(http.StatusOK, gin.H{"data": createdHook})
}

func (m HookController) ReadHookByID(c *gin.Context) {
	id := c.Param("id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal("Error parse id to int64.")
	}
	hook := hookModel.ReadHook(id64)
	if hook.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"err": "record not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": hook})
	}
}

func (m HookController) ExecuteHookByKey(c *gin.Context) {
	key := c.Param("key")
	var postData []models.OptionTemplate
	if c.Request.Method == "POST" {
		postData = []models.OptionTemplate{}
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	status := hookModel.ExecuteHookByKey(key, postData)
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (m HookController) ExecuteHookByID(c *gin.Context) {
	id := c.Param("id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal("Error parse id to int64.")
	}
	status := hookModel.ExecuteHookByID(id64, nil)
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (m HookController) UpdateHookByID(c *gin.Context) {
	var hook models.Hook
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := c.ShouldBindJSON(&hook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	hook.ID = uint(id)
	updatedHook := hookModel.UpdateHook(&hook)
	c.JSON(http.StatusOK, gin.H{"data": updatedHook})
}

// func (m HookController) UpdateHook(c *gin.Context) {
// 	var hook models.Hook
// 	if err := c.ShouldBindJSON(&hook); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}
// 	updatedHook := hookModel.UpdateHook(&hook)
// 	c.JSON(http.StatusOK, gin.H{"data": updatedHook})
// }

func (m HookController) DeleteHookByID(c *gin.Context) {
	id := c.Param("id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal("Error parse id to int64.")
	}
	status := hookModel.DeleteHook(id64)
	c.JSON(http.StatusOK, gin.H{"result": status})
}

func (m HookController) DeleteHook(c *gin.Context) {
	id := c.Param("id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal("Error parse id to int64.")
	}
	status := hookModel.DeleteHook(id64)
	c.JSON(http.StatusOK, gin.H{"result": status})
}

func (m HookController) MigrateHook(c *gin.Context) {
	hookModel.Migrate()
	c.JSON(http.StatusOK, gin.H{"result": true})
}
