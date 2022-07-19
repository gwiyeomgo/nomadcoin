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

//메세지읽기 ?
func (p *peer) read() {
	//delete peer in case of error
	//에러가 발생하면 peers 에서 peer 를 지운다
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	//& 로 pointer 만듬
	p := &peer{
		conn,
	}
	//peer 를 생성하는 순간
	go p.read()
	key := fmt.Sprintf("%s:%s", address, port)
	Peers[key] = p
	return p
}
