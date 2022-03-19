package blockchain

import (
	"errors"
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

//mempool
type mempool struct {
	Txs []*Tx
}

//mempool 곧바로 초기화 해줌 =>  비어있는 mempool
//mempool 은 memory에만 존재한다 (blockchain의 경우는 db에 저장)
var Mempool *mempool = &mempool{}

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

//transaction 을 생성해줗 makeTx

func makeTx(from, to string, amount int) (*Tx, error) {
	// gwiyeom 의 잔금이 amount 보다 금액이 적다면
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	//amount 금액과 비교했을떄 total이 작거나 같을때까지 TxIns 에 담는다.
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}
	change := total - amount
	//잔돈
	if change != 0 {
		changeTxOut := &TxOut{
			Owner:  from,
			Amount: change,
		}
		txOuts = append(txOuts, changeTxOut)
	}
	//받는사람 거래내역
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxOuts:    txOuts,
		TxIns:     txIns,
	}
	tx.getId()
	return tx, nil
}
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("gwiyeom", to, amount)
	if err != nil {
		return err
	}
	//mempool에 tx저장
	m.Txs = append(m.Txs, tx)
	return nil
}
