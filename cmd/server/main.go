package main

import (
	"flag"
	"fmt"
	service "gocopy/service/server"
	"gocopy/service/types/typesser"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/winkb/tcp1/btmsg"
)

type CopyRequest struct {
	Text string `form:"text"`
}

func handleCopy(ctx *gin.Context) {
	var b, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}
	head := &btmsg.MsgHead{
		Act: 100,
	}

	var msg = btmsg.NewMsg(head, nil)

	err = msg.FromStruct(&typesser.CopyReply{
		Text: string(b),
	})
	if err != nil {
		log.Println(err)
		return
	}

	service.GetServer().Broadcast(msg)
	fmt.Println("get data:", string(b))

}

var port int
var httpPort string

func main() {
	flag.IntVar(&port, "port", 898, "")
	flag.StringVar(&httpPort, "httpPort", "81", "")

	flag.Parse()

	go func() {
		r := gin.Default()

		r.POST("/copy", handleCopy)

		r.Run(":" + httpPort)
	}()

	service.ServerRun(port)
}
