package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"mob/websocket-api"
	"net/http"
)

var (
	upgarder = websocket.Upgrader{
		//遇到跨域问题 设置为允许
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("server listen on 0.0.0.0:7777")
	http.ListenAndServe("0.0.0.0:7777", nil)

}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		conn   *websocket_api.Connection
	)
	//升级http协议为websocket  Upgrade：websocket-api
	if wsConn, err = upgarder.Upgrade(w, r, nil); err != nil {
		return
	}
	if conn, err = websocket_api.InitConnection(wsConn); err != nil {
		goto ERR
	}
	/*for{
		//传递类型有text类型有二进制binary类型
		if msgType,data,err=conn.ReadMessage();err!=nil{
			goto ERR
		}
		fmt.Println(msgType)
		if err=conn.WriteMessage(websocket.TextMessage,data);err!=nil{
			goto ERR
		}
	}*/

	//封装后的使用方法
	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
