package controller

import (
	"net/http"

	gofstab "github.com/deniswernert/go-fstab"
	"github.com/gin-gonic/gin"
)

// GetRootPathData Get Root Path Data
func GetRootPathData(c *gin.Context) {
	mounts, err := gofstab.ParseSystem()
	if err != nil {
		errs := gin.H{
			"status":  500,
			"message": err.Error(),
		}
		c.JSON(http.StatusOK, errs)
		return
	}
	var list []string
	for _, val := range mounts {
		if val.File != "none" {
			list = append(list, val.File)
		}
	}
	data := gin.H{
		"status": 200,
		"data":   list,
	}
	c.JSON(http.StatusOK, data)
	return
}
