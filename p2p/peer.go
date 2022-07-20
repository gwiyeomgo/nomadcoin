package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// 공유를 위해
//var Peers map[string]*peer = make(map[string]*peer)

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	address string
	key     string
	port    string
}

//recevier 함수는 변형시킬 필요가 있을 때 사용한다

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()
	var keys []string

	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

//메세지읽기 ?
func (p *peer) read() {
	//delete peer in case of error
	//에러가 발생하면 peers 에서 peer 를 지운다
	defer p.close()
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

//err 있거나 채널이 닫혔을 때
func (p *peer) close() {
	// peers struct(Peers) 변수는
	//누구에게나 잠겨있는 상태가 된다
	//만약 Peers 변수에 접근하려는 또 다른 go 루틴이 있어도
	//우리가 잠금을 풀어 줄 때까지 기다려야 한다
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) write() {
	//defer 는 함수가 종료된 후에 어떠한 코드를 실행시킬 수 있다
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	//& 로 pointer 만듬
	p := &peer{
		conn,
		make(chan []byte),
		address,
		key,
		port,
	}
	//peer 를 생성하는 순간
	go p.read()
	go p.write()

	Peers.v[key] = p
	return p
}
