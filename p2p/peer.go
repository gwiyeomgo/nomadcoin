package p2p

import "github.com/gorilla/websocket"

// 공유를 위해
var Peers

type peer struct {
	conn *websocket.Conn
}