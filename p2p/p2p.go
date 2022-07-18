package p2p

import (
	"github.com/gorilla/websocket"
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
	//해당 function 에서는 upgrade 역할만 한다

	//equest origin not allowed by Upgrader.CheckOrigin
	//아무나 너의 서버에 접속할 수 있게 하면 안되기 때문에... 에러 발생
	//CheckOrigin 은 유요한 webSocket 연결인지
	// authenticate 인증할 때 사용
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	//websocket pacakge 가 있지만 일부기능 동작 X
	//https://github.com/gorilla/websocket 사용
	//go get github.com/gorilla/websocket

	//1. API 에서 request 를 가져온다
	//2. upgrader.Upgrade()
	//http connection 을 WS connection 으로 upgreade 했음
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
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
