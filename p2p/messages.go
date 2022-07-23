package p2p

import (
	"encoding/json"
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/blockchain"
	"github.com/gwiyeomgo/nomadcoin/utils"
)

type MessageKind int

/*const (
	MessageNewestBlock       MessageKind = 1 // 새로운 블록을 얻었을 때
	MessageAllBlocksRequest  MessageKind = 2 //모든 블록을 요청
	MessageAllBlocksResponse MessageKind = 3  // 모든 블록을 응답
)*/
//iota 라는 타입
//dl type 과 value 를 기본적으로 자동 생성
//단 첫번째 변수를 정의하고 변수의 타입을 정해야 한다

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

//Payload
//한 개의 block,peer address,transaction 일 수 있다
type Message struct {
	Kind    MessageKind
	Payload []byte
}

/*
//Payload 를 json 인코딩 후 추가
func (m *Message) addPayload(p interface{}) {
	//interface => json 바이트
	//Marshal 은 json 인코딩 된 b 를 반환
	b,err := json.Marshal(p)
	utils.HandleErr(err)
}*/

//메세지를 구성하는 함수
func makeMessage(kind MessageKind, payload interface{}) []byte {
	//payload interface{} =>  인터페이스는 무엇이든 가져올꺼고 ,바이트로 반환한다
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	//메세지를 JSON 으로 인코딩
	/*mJson, err := json.Marshal(m)
	utils.HandleErr(err)*/
	return utils.ToJSON(m)
}

//메세지 => json 으로
func sendNewestBlock(p *peer) {
	fmt.Printf("sending newest block to  %s\n", p.key)

	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	//가장 최신의 block 을 찾고
	//우리가 보낼 메세지 type 와 Payload 를 가져
	m := makeMessage(MessageNewestBlock, b)
	//메세지가 준비되면 그 메세지를 inbox 로 보낸다
	p.inbox <- m
}

//모든 블록을 요청한다
func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}
func sendAllBlocks(p *peer) {
	blocks := blockchain.Blocks(blockchain.Blockchain())
	m := makeMessage(MessageAllBlocksResponse, blocks)
	p.inbox <- m
}

//메세지를 처리해 주는 함수
func handleMsg(m *Message, p *peer) {
	//fmt.Printf("Peer:%s,Sent a message with kind of:%d",p.key,m.Kind)
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)

		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		//fmt.Println(payload)
		// 현재 블록은?
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		//두 database 가 하나의 block 만 가지고 있기 때문에 >=
		if payload.Height >= b.Height {
			//payload.Height > our block
			// request all the blocks from 4000
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			// send 4000 our block
			fmt.Printf("Sending newest blocks to %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		//모든 블록을 달라고 요청하면 전부 보낸다
		fmt.Printf("%s wants all the blocks\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().Replace(payload)
	}
}
