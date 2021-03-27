package controller

import (
	"net/http"
	"pulltg/utils"

	"github.com/gin-gonic/gin"
)

// GetConfigData Get Config Data
func GetConfigData(c *gin.Context) {
	runPath, _ := c.Get("runPath")
	config := utils.GetConfig(runPath.(string))
	data := gin.H{
		"status": 200,
		"config": config.PathList,
	}
	c.JSON(http.StatusOK, data)
}
