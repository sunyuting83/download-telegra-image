package controller

import (
	"fmt"
	"net/http"
	"strings"

	gofstab "github.com/deniswernert/go-fstab"
	"github.com/gin-gonic/gin"
)

// Download download
func Download(c *gin.Context) {
	var x bool = true
	var w string = c.DefaultQuery("w", "1")
	var u string = c.Query("url")
	if w != "1" {
		x = false
	}
	mounts, _ := gofstab.ParseSystem()

	for _, val := range mounts {
		fmt.Printf("%v\n", val.File)
	}
	if len(u) <= 0 {
		error := gin.H{
			"status":  500,
			"message": "URL cannot be empty",
		}
		c.JSON(http.StatusOK, error)
		return
	}
	if !strings.Contains(u, "https://telegra.ph") {
		error := gin.H{
			"status":  500,
			"message": "URL must be telegra",
		}
		c.JSON(http.StatusOK, error)
		return
	}
	data, title, b := Scrape(x, u)
	datas := gin.H{
		"status":  200,
		"data":    data,
		"title":   title,
		"statuss": b,
	}
	c.JSON(http.StatusOK, datas)
	return
}
