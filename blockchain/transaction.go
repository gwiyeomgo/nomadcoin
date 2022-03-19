package blockchain

import (
	"github.com/gwiyeomgo/nomadcoin/utils"
	"time"
)

const (
	minerReward int = 50
)

//거래
type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

//Transacion을  hash값으로 바꾸고 Id 값으로 넣음
func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

//입력갑
type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

//출력값
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

//채굴자를 주소로 삼는 코인베이 거래내역을 생성해서 Tx포인터를 retrun
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"CoinBase", minerReward},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// 어떤 사용자 ,혹은 주소가 블록체인에 자산을 얼마나 갖고 있는지 찾아내는 함수
// 채굴자 주소가 소유중인 모든 출력값을 찾아라
