package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// 공유를 위해
var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn *websocket.Conn
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	//& 로 pointer 만듬
	p := &peer{
		conn,
	}
	key := fmt.Sprintf("%s:%s", address, port)
	Peers[key] = p
	return p
}
