package gonet

import (
	"base/glog"
	"net/http"

	"github.com/gorilla/websocket"
)

type IWebSocketServer interface {
	OnWebAccept(conn *websocket.Conn)
}

type WebSocketServer struct {
	WebDerived IWebSocketServer
}

var upgrader = websocket.Upgrader{} // use default options

func (this *WebSocketServer) WebBind(addr string) error {
	http.HandleFunc("/", this.WebListen)
	err := http.ListenAndServe(addr, nil)
	if nil != err {
		glog.Error("[WebSocketServer] init failed", addr)
	}
	return nil
}

func (this *WebSocketServer) WebListen(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		// allow all connections by default
		return true
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error("[WebSocketServer] failed to upgrade", r.RemoteAddr, err)
		return
	}

	glog.Info("[WebSocketServer] recv connect ", r.RemoteAddr, w.Header().Get("Origin"), r.Header.Get("Origin"))

	this.WebDerived.OnWebAccept(c)
}

func (this *WebSocketServer) WebClose() error {
	return nil
}
