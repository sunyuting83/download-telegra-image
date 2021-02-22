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
	api.Use(utils.CORSMiddleware())
	{
		api.GET("/download", utils.SetConfigMiddleWare(d), controller.Download)
	}

	return router
}
