package websocket_api

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	isClose   bool
	m         sync.Mutex
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	go conn.readLoop()
	go conn.writeLoop()
	return
}

//读取信息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

//信息放入通道
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connect is close")

	}
	return
}

func (conn *Connection) Close() {
	//线程安全
	conn.wsConn.Close()
	conn.m.Lock()
	if !conn.isClose {
		close(conn.closeChan)
	}
	conn.m.Unlock()
}

//从信息通道中读取信息
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto Err
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto Err
		}

	}
Err:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto Err
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto Err
		}
	}
Err:
	conn.Close()
}
