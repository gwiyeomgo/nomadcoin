package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gwiyeomgo/nomadcoin/blockchain"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"net/http"
)

/*
//webSocket 은 alive connection 이다
//ex - wife 는 bi-directiional(양 방향) connection 이고 state 가 full 이다
/*
1.서버로 request 를 보내고
2.서버에게 지금 연결을 webSocket 연결로 upgrade 하지고 한다
3. webSocket 연결을 즉시 보내서, upgrader 를 go 로 만듬
2개의 프로토콜 이 있는데 http,ws 이다
http 는 stateless 이고 ws 는 stateful 이다


*/
var conns []*websocket.Conn
var upgrader = websocket.Upgrader{}

func Upgreade(rw http.ResponseWriter, r *http.Request) {
	//port 3000d이 port 4000 에서 온 request 를 will be upgrade
	//해당 function 에서는 upgrade 역할만 한다

	// 원격 주소를 사용하면 서버가 요처을 보낸 네트워크 주소를 기억할 수 있다.
	// RemoteAddr allows HTTP servers and other software to record
	// the network address that sent the request,
	//fmt.Println(r.RemoteAddr)
	//문제는 127.0.0.1:62279 원격주소는 저장되지만..
	//:62279 가 아닌 peer 에서 열려 있는 포트를 저장하고 싶다
	// r.RemoteAddr 를 :  로 쪼갠다
	//result := strings.Split(r.RemoteAddr,":")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	//initPeer(conn,result[0],"xx")
	//openPort
	//request query 받아오기
	openPort := r.URL.Query().Get("openPort")
	//openPort 가 "" 인 경우

	//equest origin not allowed by Upgrader.CheckOrigin
	//아무나 너의 서버에 접속할 수 있게 하면 안되기 때문에... 에러 발생
	//CheckOrigin 은 유요한 webSocket 연결인지
	// authentuicate 인증할 때 사용
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}

	fmt.Printf("%s want an upgrade \n", openPort)

	//websocket pacakge 가 있지만 일부기능 동작 X
	//https://github.com/gorilla/websocket 사용
	//go get github.com/gorilla/websocket

	//1. API 에서 request 를 가져온다
	//2. upgrader.Upgrade()
	//http connection 을 WS connection 으로 upgreade 했음
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	//peer 들에게 메세지를 보내는 방법?
	// 읽는 방법 ? conn.ReadMessage()

	// conn 은 port 3000  과 4000 을 이여준다
	//adreess 와 port 정보 어떻게 알지?
	//initPeer(conn,"xx","xx")

	// 누구와 연결되어 있고 그들의 열린 포트가 무엇인지 알 수 있다
	//앞으로는 우리와 연결된 peers 목록과
	//열린 포트가 무엇인지 저장해야 한다
	// IP 리스트와 peers 의 오픈포트를 저장하는 이유는
	//새 node 가 들어오면,
	//ex 2000 포트가 들어오면 3000 에 연결하고
	//3000 번 포트가 연결되어있는 모든 peers 를 보게 될거고
	//3000번 포트가 4000번 포트에 업그레이드 요청을 보낼 수 있따
	initPeer(conn, ip, openPort)

	//time.Sleep(20 * time.Second)
	//	conn.WriteMessage(websocket.TextMessage, []byte("Hello form Port 3000!"))
	//peer.inbox <- []byte("Hello form Port 3000!")

	/*	fmt.Println("Waiting 4 message...")
		//3. coonn 으로 WriteMessage, WriteJSON  또는 ReadMessage
		_, p ,err := conn.ReadMessage()
		fmt.Println("Message arrived")
		utils.HandleErr(err)
		//메세지를 읽음
		fmt.Printf("%s",p)

		//한번 더, message 를 읽는 것은 block 이다
		// message 가 도착하면, unblock 하고 받고 print 하지
		fmt.Println("Waiting 4 message...")
		_, p ,err = conn.ReadMessage()
		fmt.Println("Message arrived2")
		utils.HandleErr(err)
		fmt.Printf("%s",p)*/

	//_, p ,err := conn.ReadMessage() 이 block 하기 때문에
	// for 문 밑에 conn.WriteMessge() 해도 실행 X
	// 2개 function 실행을 위해서 go 루틴 사영영

	/* 같은 브라우저에서 메세지 공유
	for {
		//fmt.Println("Waiting message...")
		_, p, err := conn.ReadMessage()
		//fmt.Println("Message arrived...")
		//새로고침하면 err 발생
		//utils.HandleErr(err) 는 panic 하니까 err 있다면
		if err != nil {
			//err 가 있다면 연결 끊음
			//conn.Close()
			break
		}
		fmt.Printf("Just got : %s\n\n", p)
		//매번 메세지를 받을때마다 이런 내용을 출력
		time.Sleep(5 * time.Second)
		//5초기다리고.유저에게 새 메세지 보냄
		message := fmt.Sprintf("New message: %s", p)
		utils.HandleErr(conn.WriteMessage(websocket.TextMessage, []byte(message)))

	}*/
	//Upgrade 는 go 루틴으로 돌아감 => 덕북에 다른 브라우저로 전송 가능

	//다른 브라우저로 메세지 전송송
	//연결들을 배열에  추가
	//conns = append(conns,conn)
	/*
		for {
			//* 중요 : read 할 때 block 했고, 메세지 보냄
			_,p, err := conn.ReadMessage()
			if err != nil {
				break
			}
			//브라우저가 다를 때 메세지 전송하도록 코딩
			//문제는 새로고침하면 connection 이 끊긴다
			//conns 에 여전히 연결이 남아있기 때문에
			for _, aConn := range conns {
				//aConn 이 닫혔는지 확인?
				if aConn != conn {
					utils.HandleErr(aConn.WriteMessage(websocket.TextMessage, p))
				}
			}
		}*/
	//읽기 , 쓰기 용도에 따라 분리
	//특정 메세지를 1개 특정 브라우저로 보내고 싶다

	//conn.WriteMessage(websocket.TextMessage)
}

func AddPeer(address, port, openPort string, broadcast bool) {
	//bool 추가해
	//이게 첫 연결인지 아니면,broadcast 를 통한 연결인지 확인

	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	//go에서 connection 하기
	//이 URL 을 call 하면 새로운 connection 을 만든다
	//websocket 서버랑 연결하려고 할 때 upgrade 하기전 앞으로 도착할 request 를 체크할 수 있다
	//만약 websocket 으로 인증하고 싶다면 requestHeader 에 token 을 보낸다

	// 여기는 port:4000 이고 port:3000 으로 연결할길원함 (dial)
	//연결을 시도할 때 우리가 연결하려는 서버에게
	//어떤 포트가 열려있는지도 알려주자
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	//time.Sleep(10 * time.Second)
	//conn.WriteMessage(websocket.TextMessage, []byte("Hello form Port 4000!"))
	//peer.inbox <- []byte("Hello form Port 4000!")
	if broadcast {
		//broadcast 는 새로운 peer 가 rest api 를 통해서 올 때만 true 가 된다
		broadcastNewPeer(peer)
		return
	}

	sendNewestBlock(peer)

}

func notifyNewestBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m

}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewestBlock(b, p)
	}
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	//makeMessage 는 메세지를 만든 다음 payload 와 메세지 자체를 json 으로 변환
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m

}
func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}
func notifyNewPeer(payload string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, payload)
	p.inbox <- m
}
func broadcastNewPeer(newPeer *peer) {
	//업그레이드 받고 해당 함수를 실행
	// 새로 들어온 peer 에게는 메시지를 보내지 않고 나머지 peer 에 보냄
	for key, p := range Peers.v {
		if newPeer.key != key {
			//보낸사람이 누구인지 알고 싶음
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
