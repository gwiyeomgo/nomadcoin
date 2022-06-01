package blockchain

import (
	"errors"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"github.com/gwiyeomgo/nomadcoin/wallet"
	"time"
)

const (
	minerReward int = 50
)

//거래
type Tx struct {
	ID        string   `json:"id"`
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

func validate(tx *Tx) bool {
	//transaction input 에 참조된
	//transaction output 을 소유한 사람을 검증
	//transaction을 만드려면 output 이 필요하다
	//그 output은 다음 transactin을 만들 떄 input 이 된다
	//근데 우린 우리가 그 output 을 소유하고 있다는 것을 증명해야 한다
	valid := true
	for _, txIn := range tx.TxIns {
		//tansaction id 로 이전 트랜잭션 알 수 있따
		preTx := FindTx(Blockchain(), txIn.TxID)
		//이전 transaction 이 blockchain 에 없다면
		if preTx == nil {
			valid = false
			break
		}
		//유효하다면 address = 이 transaction input 이 참조한
		// transaction output 의 주소
		address := preTx.TxOuts[txIn.Index].Address
		//address = publick key 로 서명을 검증할 수있다
		// 우리가 검증하고자 하는 것 = payload = transactin ID
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}

	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
	//2. label 사용
	//여러 개의 for loop 이 중첩됐을 경우
	//바깥쪽 for loop을 종료시킬 방법 label
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			//아래 코드는 멈추지 않음
			//uOut 과 같은 트랜잭션 ID와 index 를 가지고 이는 항목이 있는지 확인
			if uTxOut.TxID == input.TxID && input.Index == uTxOut.Index {
				exists = true
				//break //해당 for 문만 멈춤춤				//1.원하는 값을 찾으면, true 반환시켜 함수를 끝낸다.
				break Outer
				//return true
			}
		}
	}
	return exists
}

//Transacion을  hash값으로 바꾸고 Id 값으로 넣음
func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	//transaction 의 모든 transaction input 들에 서명을 저장한다
	for _, txIn := range t.TxIns {
		//우리가 갖고있는 Wallet 의 private key 로
		//transactin id 에 서명한다
		//그 서명을,우리가 서명한 id 를 갖는 트랜잭션의 transaction input 에 저장했다다
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet())
	}
}

//입력갑
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"` //TxIn의 ID 와 Index 는 transactin 이 이 input 을 생성한 output을 가지고 있는지 알려준다.
	Signature string `json:"signature"`
	//Owner string `json:"owner"`
	//	Amount int    `json:"amount"`
}

//출력값
// address 가 바로 사람들이 너에게 코인을 보내는 곳
type TxOut struct {
	Address string `json:"address"` //public key 를 string 으로 한거
	//Owner  string `json:"owner"`
	Amount int `json:"amount"`
}

//사용하지 않은 outs
type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

//채굴자를 주소로 삼는 코인베이 거래내역을 생성해서 Tx포인터를 retrun
func makeCoinbaseTx(address string) *Tx {
	/*txIns := []*TxIn{
		{"CoinBase", minerReward},
	}*/
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

//transaction 을 생성해줗 makeTx

var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("Tx Invalid")

func makeTx(from, to string, amount int) (*Tx, error) {
	// gwiyeom 의 잔금이 amount 보다 금액이 적다면
	if BalanceByAddress(Blockchain(), from) < amount {
		return nil, ErrorNoMoney
	}
	//amount 금액과 비교했을떄 total이 작거나 같을때까지 TxIns 에 담는다.
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	//보내는 사람의 소비안된 거내내역
	uTxOuts := UTxOutsByAddress(Blockchain(), from)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{
			uTxOut.TxID,
			uTxOut.Index,
			from,
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	/*oldTxOuts := Blockchain().TxOutsByAddress(from)
	//현재 상태는 transaction output 의 유일성을 확인하지 않고 있다.
	// 보내는이의 계정에 있는 돈보다 많은 돈을 보내려는겋 확인하고 있지 않다
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}
	*/
	//잔돈 계산
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	//받는사람 거래내역
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxOuts:    txOuts,
		TxIns:     txIns,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	//mempool에 tx저장
	m.Txs = append(m.Txs, tx)
	return nil
}

//승인할 트랜잭션들을 가져오기
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	// mempool 에서 transaction 비워주기
	m.Txs = nil
	return txs
}
