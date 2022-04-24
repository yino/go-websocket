package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/yino/websocket/client"
)

var Client []*client.Client

// main .
func main() {
	r := gin.Default()
	r.GET("/ws", ws)

	r.Run("0.0.0.0:8080")
}

func ws(ctx *gin.Context) {
	up := websocket.Upgrader{
		// 设置允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := up.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("conn error", err)
		return
	}
	client := client.NewClient(conn)
	Client = append(Client, client)
}
