package main

import (
	"flag"
	"pulltg/router"
	"pulltg/ws"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		port string
		d    string
	)
	flag.StringVar(&port, "p", "13002", "default port")
	flag.StringVar(&d, "d", "/etc/roubian", "default path")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := router.InitRouter(d, port)
	go ws.Manager.Start()
	router.Run(strings.Join([]string{":", port}, ""))
}
