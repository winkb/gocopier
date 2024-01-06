package main

import (
	"flag"
	"fmt"
	"gocopy/lib/scopy"
	service "gocopy/service/client"
	"gocopy/service/types/typesser"

	"github.com/winkb/tcp1/btmsg"
)

var address = ""

func main() {
	flag.StringVar(&address, "host", "127.0.0.1:898", "")
	flag.Parse()

	service.ClientRun(address, func(msg btmsg.IMsg, req typesser.CopyReply) {
		fmt.Println("receive:\n" + req.Text)
		scopy.CopyText("复制：\n" + req.Text)
	})
}
