package service

import (
	"fmt"
	"gocopy/service/types/typesser"

	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/mytcp"
)

type RouteHandle func(msg btmsg.IMsg)

type RouteInfo struct {
	Handle RouteHandle
}

func newMsg(act uint16, req any) btmsg.IMsg {
	hd := btmsg.NewMsgHead()
	hd.Act = act
	res := btmsg.NewMsg(hd, nil)
	err := res.FromStruct(req)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

var routes = map[uint16]*RouteInfo{}

func init() {
	routes[100] = &RouteInfo{
		Handle: func(msg btmsg.IMsg) {
			handleCopy(msg, parseReq(typesser.CopyReply{}, msg))
		},
	}
}

func parseReq[T any](v T, msg btmsg.IMsg) T {
	_, _ = msg.ToStruct(&v)
	return v
}

var myCopyHandle func(msg btmsg.IMsg, req typesser.CopyReply)

func handleCopy(msg btmsg.IMsg, req typesser.CopyReply) {
	myCopyHandle(msg, req)
}

func ClientRun(address string, copyHandle func(msg btmsg.IMsg, req typesser.CopyReply)) {
	cli := mytcp.NewTcpClient(address)

	myCopyHandle = copyHandle

	cli.OnReceive(func(v btmsg.IMsg) {
		act := v.GetAct()
		r, ok := routes[act]
		if !ok {
			fmt.Println("not found handle", act)
			return
		}

		r.Handle(v)
	})

	cli.OnClose(func(isServer bool, isClient bool) {
		if isClient {
			fmt.Println("服务端断开连接")
		}

		if isServer {
			cli.ReleaseChan()
			fmt.Println("我自己端口连接")
		}
	})

	wg, err := cli.Start()
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
