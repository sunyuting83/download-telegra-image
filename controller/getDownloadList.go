package controller

import (
	"net/http"
	"pulltg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetDownloadList Get Download List
func GetDownloadList(c *gin.Context) {
	runPath, _ := c.Get("runPath")
	config := utils.GetConfig(runPath.(string))
	dataFileName := strings.Join([]string{config.RunPath, config.DataFile}, "/")
	list := GetDataFile(dataFileName)
	data := gin.H{
		"status": 200,
		"list":   list,
	}
	c.JSON(http.StatusOK, data)
	return
}
