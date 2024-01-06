package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/mytcp"
	"github.com/winkb/tcp1/util/numfn"
)

type RouteHandle func(conn *mytcp.TcpConn, msg btmsg.IMsg)

type RouteInfo struct {
	Handle RouteHandle
	Info   any
}

var server mytcp.ITcpServer

type CopyToClient struct {
	Text string
}

var routes = map[uint16]*RouteInfo{}

func GetServer() mytcp.ITcpServer {
	return server
}

func init() {
	routes[0] = &RouteInfo{
		Handle: func(conn *mytcp.TcpConn, msg btmsg.IMsg) {
			HandleDefault(conn, msg, nil)
		},
	}
}

func parseReq[T any](v T, msg btmsg.IMsg) T {
	_ = msg.FromStruct(&v)
	return v
}

func logHandle(name string, t time.Time) func() {
	return func() {
		fmt.Println("handle", name, "in")
		fmt.Println("handle", name, "out", numfn.ToStr(time.Now().Sub(t).Nanoseconds())+"ns")
	}
}

func HandleDefault(conn *mytcp.TcpConn, msg btmsg.IMsg, req any) {
	defer logHandle("default", time.Now())

	fmt.Println("sever receive default msg ", req)
}

func ServerRun(port int) {
	server = mytcp.NewTcpServer(fmt.Sprintf("%d", port), btmsg.NewReader())
	wg, err := server.Start()
	if err != nil {
		panic(err)
	}

	server.OnClose(func(conn *mytcp.TcpConn, isServer bool, isClient bool) {
		if isClient {
			fmt.Println("客户端断开连接")
		}

		if isServer {
			fmt.Println("我自己断开连接")
		}
	})

	server.OnReceive(func(conn *mytcp.TcpConn, msg btmsg.IMsg) {
		act := msg.GetAct()
		hv, ok := routes[act]
		if !ok {
			fmt.Println("not found handle", act)

			// 走默认路由
			act = 0
			hv = routes[act]
		}

		hv.Handle(conn, msg)
	})

	chSingle := make(chan os.Signal)

	signal.Notify(chSingle, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		v := <-chSingle
		switch v {
		case syscall.SIGINT:
			fmt.Println("ctr+c")
		case syscall.SIGTERM:
			fmt.Println("terminated")
		}

		server.Shutdown()

		fmt.Println(v)
	}()

	wg.Wait()
	close(chSingle)
}
