package main

import (
	"flag"
	"pulltg/router"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		port string
	)
	flag.StringVar(&port, "p", "13002", "default port")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := router.InitRouter()
	router.Run(strings.Join([]string{":", port}, ""))
}
