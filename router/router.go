package router

import (
	controller "pulltg/controller"
	utils "pulltg/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(d string) *gin.Engine {
	// router := gin.New()
	router := gin.Default()
	api := router.Group("/api")
	api.Use(utils.CORSMiddleware(), utils.SetConfigMiddleWare(d))
	{
		api.GET("/download", controller.Download)
		api.GET("/getconfig", controller.GetConfigData)
		api.GET("/getpath", controller.GetPathData)
		api.GET("/getrootpath", controller.GetRootPathData)
		api.PUT("/setconfig", controller.SetConfigPathData)
		api.GET("/downlist", controller.WsPage)
	}

	return router
}
