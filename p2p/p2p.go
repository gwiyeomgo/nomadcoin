package p2p

import (
	"github.com/gorilla/websocket"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"net/http"
)

//gorilla
//go get github.com/gorilla/websocket

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	//웹소켓에 아무나 접근하면 안돼니까~
	//upgrade 할 때 이게 필요하다=> 허가된 사람만 접속시킨다 => request를 받아서 origin cookie 등 확인
	//webSoket 연결이 허가됨
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	//connection 을 만듬
	//http connection 을 WS connection 으로 upgrade 했음
	_, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
}
