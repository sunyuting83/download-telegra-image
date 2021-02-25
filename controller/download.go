package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"pulltg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// Download download
func Download(c *gin.Context) {
	var (
		x    bool   = true
		w    string = c.DefaultQuery("w", "1")
		u    string = c.Query("url")
		p    string = c.Query("path")
		errs gin.H  = GetErrorMessage("some error")
	)
	downPath, _ := url.QueryUnescape(utils.DecodeBytes(p))

	if w != "1" {
		x = false
	}
	runPath, _ := c.Get("runPath")
	config := utils.GetConfig(runPath.(string))
	if len(config.PathList) <= 0 {
		errs = GetErrorMessage("Must set download path")
		c.JSON(http.StatusOK, errs)
		return
	}
	if len(u) <= 0 {
		errs = GetErrorMessage("URL cannot be empty")
		c.JSON(http.StatusOK, errs)
		return
	}
	if len(p) <= 0 {
		errs = GetErrorMessage("Download Path cannot be empty")
		c.JSON(http.StatusOK, errs)
		return
	}
	if !strings.Contains(u, config.RootURL) {
		errs = GetErrorMessage("URL must be telegra")
		c.JSON(http.StatusOK, errs)
		return
	}
	data, title, b := Scrape(x, u, config.Cors, config.RootURL)
	DownloadPath := strings.Join([]string{downPath, title}, "/")
	checkPath, err := utils.PathExists(DownloadPath)
	if err != nil {
		errs = GetErrorMessage(err.Error())
		c.JSON(http.StatusOK, errs)
		return
	}
	if checkPath {
		errs = GetErrorMessage("The directory already exists")
		c.JSON(http.StatusOK, errs)
		return
	}
	if !b {
		errs = GetErrorMessage("Failed to get data")
		c.JSON(http.StatusOK, errs)
		return
	}
	// save data to file
	dataFileName := strings.Join([]string{config.RunPath, config.DataFile}, "/")
	// fmt.Println(utils.MakeMD5(DownloadPath))
	df := utils.GetDataFile(dataFileName)
	df.Running = append(df.Running, utils.SaveData{Total: len(data), Completed: 0, Key: utils.MakeMD5(DownloadPath), Path: DownloadPath})
	saveData, _ := json.Marshal(df)
	_ = ioutil.WriteFile(dataFileName, saveData, 0644)
	port, _ := c.Get("port")
	// make download path
	os.MkdirAll(DownloadPath, os.ModePerm)
	// start download
	go DownloadImages(data, DownloadPath, dataFileName, port.(string))
	datas := gin.H{
		"status": 200,
		"data":   data,
		"title":  title,
	}
	c.JSON(http.StatusOK, datas)
	return
}

// GetErrorMessage get error message
func GetErrorMessage(e string) gin.H {
	return gin.H{
		"status":  500,
		"message": e,
	}
}
