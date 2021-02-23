package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pulltg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// Path path
type Path struct {
	PATH string `form:"path" json:"path" xml:"path"  binding:"required"`
}

// SetConfigPathData Set Config Path Data
func SetConfigPathData(c *gin.Context) {
	var form Path
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.PATH) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't path",
		})
		return
	}

	runPath, _ := c.Get("runPath")
	config := utils.GetConfig(runPath.(string))
	if checkHasPath(form.PATH, config) {
		data := gin.H{
			"status":  200,
			"message": "has data",
		}
		c.JSON(http.StatusOK, data)
		return
	}
	titleSplit := strings.Split(form.PATH, "/")
	l := len(titleSplit) - 1
	title := titleSplit[l]
	config.PathList = append(config.PathList, utils.Path{Title: title, Path: form.PATH})
	jsonFile := strings.Join([]string{runPath.(string), "config.json"}, "/")
	saveData, _ := json.Marshal(config)
	_ = ioutil.WriteFile(jsonFile, saveData, 0755)
	data := gin.H{
		"status":  200,
		"message": "OK",
	}
	c.JSON(http.StatusOK, data)
	return
}

// checkHasPath check has path
func checkHasPath(p string, list *utils.Config) bool {
	var x bool = false
	for _, item := range list.PathList {
		if p == item.Path {
			x = true
			break
		}
	}
	return x
}
