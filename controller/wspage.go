package controller

import (
	"net/http"
	"pulltg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// WsPage is a websocket handler
func WsPage(c *gin.Context) {
	// change the reqest to websocket model
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	runPath, _ := c.Get("runPath")
	// websocket connect
	client := &ws.Client{ID: uuid.NewV4().String(), Socket: conn, Send: make(chan []byte), RootPath: runPath.(string)}

	ws.Manager.Register <- client

	go client.Read()
	go client.Write()
}
