package controller

import (
	"net/http"
	"net/url"
	"path/filepath"
	"pulltg/utils"

	"github.com/gin-gonic/gin"
)

// GetPathData Get Path Data
func GetPathData(c *gin.Context) {
	var (
		p string = c.Query("path")
	)
	path, _ := url.QueryUnescape(utils.DecodeBytes(p))
	pathData, err := GetPathName(path)
	if err != nil {
		errs := gin.H{
			"status":  500,
			"message": err.Error(),
		}
		c.JSON(http.StatusOK, errs)
		return
	}
	data := gin.H{
		"status": 200,
		"data":   pathData,
	}
	c.JSON(http.StatusOK, data)
	return
}

// GetPathName Get Path Name
func GetPathName(p string) (path []string, err error) {
	filepathNames, err := filepath.Glob(filepath.Join(p, "*"))
	if err != nil {
		return path, err
	}

	for i := range filepathNames {
		if utils.IsDir(filepathNames[i]) {
			path = append(path, filepathNames[i])
		}
	}
	return path, nil
}
