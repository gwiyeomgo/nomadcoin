package blockchain

import (
	"github.com/gwiyeomgo/nomadcoin/db"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"sync"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newest_hash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

/*
//bytes decode
//db에서 찾은 byte를 텅빈 블록체인의 memory address 에 decode
func (b *blockchain) fromBytes(data []byte){
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
	//utils.HandleErr(gob.NewDecoder(bytes.NewReader(data)).Decode(b))
}*/

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

//블록이 새로 만들어 질떄마다 블록과 블록체인의 상황을 db에 저장
//db에 block 을 저장하는 코드
//블록체인을 처음 실행시키는 사람의 관점
//func (b *blockchain) AddBlock(data string) {
func (b *blockchain) AddBlock() {
	//block := createBlock(data, b.NewestHash, b.Height+1)
	//아래 코드처럼 수정하면 deadlock 발생 X
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	//difficulty 지정
	//block 추가할때마다 chain에 difficulty 변경됨
	b.CurrentDifficulty = block.Difficulty
	persist(b)
}

//block에 checkoutpoint 를 지정
//func (b *blockchain) persist() {
//method 가 아닌 함수로
func persist(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

//예상 시간은 5개의 블록이 2분마다 생성되는 시간 => 10분
//10분동안 5개 블록 생성되길 예상
//blockchain을 변화시키지 않음
//func (b *blockchain) recalculateDifficulty() int {
func recalculateDifficulty(b *blockchain) int {
	//모든 blocks 를 받아온다
	allBlocks := Blocks(b)
	//최초 block 0 번째
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	//두 블록들의 생성 사이에 걸린 시간
	//Timestamp 는 분단위로 초단위 값이 필요하기 때문에 60 나눔
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	//예상시간
	expectedTime := difficultyInterval * blockInterval
	// 실제 예상 시간보다 적다면 (너무 빨리 생성) -> difficulty 증가
	//실제 시간이 예를 들어 9분이면
	//10분에 근접했으니 difficulty 를 재설정하지 않음
	//만약 실제시간이 11분이라면 큰자이가 아니기때문에 difficulty 줄이지 않음
	// allowedRange 2분차이는 는접
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		//너무 어렵
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

//difficulty 함수를 blockchain 패키지 내에서 사용가능
//block.go 에 상수를 지움
//func (b *blockchain) difficulty() int {
func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//만약 나머지값이 0이라면
		// recalculate the difficulty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

//b가 nil 인 경우는 딱 한번만 발생
func Blockchain() *blockchain {
	//if b == nil { //이부분 지우면 DeadLock 발생
	//b를 초기화 (처음이자 마지막)
	//Do 함수는 매개별수로 온 function이 끝날때까지 멈추지 않아 DeadLock 바랭
	//함수가 자기 자신을 영원히 호출
	once.Do(func() {
		//우리가 몇 천 개의 goroutine을 실행해도 오직 한번만 호출

		b = &blockchain{
			Height: 0,
		}
		// search for checkpoint on the db (db에서 블록체인을 가져오는 함수)
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			//b.AddBlock("Genesis")
			b.AddBlock()
		} else {
			//fmt.Println("...Restore")
			//restore b from bytes (db 에는 bytes로 저장되어 있음)
			b.restore(checkpoint)
		}
	})
	//}
	//	fmt.Printf("NewesHash:%s\nHeight:%d\n", b.NewestHash, b.Height)
	return b
}

//

//모든 block 을 보여주는 함수
//Block 포인터의 slice 를 retrun
//func (b *blockchain) Blocks() []*Block {
//blockchain을 변화시키지 않고 그저 input 으로만 사용
//메서드가 아닌 function 으로 변경
func Blocks(b *blockchain) []*Block {
	//NewestHash 를 갖고 해당 블록을 찾는다.
	//prevHash 가 없는 블록을 찾을 때 까지
	//처음에는 블록체인의 NewestHash 찾음
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PreHash != "" {
			hashCursor = block.PreHash
		} else {
			break
		}
	}
	return blocks
}

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	//blockchain 에 있는 모든 block 들에 대해서
	//이 transaction 들을 txs 에 넣어줬다다
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transaction...)
	}
	return txs
}

// transaction id 로
func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

//전체 outs 는 필요 없고 unspnt trnasacion output 만 필요 아래 코드 지움
// 어떤 사용자 ,혹은 주소가 블록체인에 자산을 얼마나 갖고 있는지 찾아내는 함수
// 채굴자 주소가 소유중인 모든 출력값을 찾아라
//1. 모든 거래 출력값만 반환하는 함수
/*func (b *blockchain) txOuts() (txOuts []*TxOut) {
	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transaction {
			//배열을 합할때 ... 사용
			//모든 거래내역의 출력값을 하나의 슬라이스에 모음
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

//2.거래 출력값들을 주소에 따라 걸러 낸다
func (b blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
}
*/
//우린 input 으로 ㅏ용된 output 을 소유한 트랜잭션들을 마킹했다 ID로
// 그래서 만약 우리가 지금 그 트랜잭션 내부에 있다면
// 우린 그 트랜잭션의 output 으로 upsent transaction output 을 생성하지 않을거야
//input 에서 사용되지 않은 output 을
//func (b blockchain) UTxOutsByAddress(address string) []*UTxOut {
func UTxOutsByAddress(b *blockchain, address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	//1.전체 블록을 가져와 각각의 transaction 의 input 에서
	//adress 가 같은 txIds 를 뽑는다.
	for _, block := range Blocks(b) {
		for _, tx := range block.Transaction {
			for _, input := range tx.TxIns {
				//if input.Owner == address {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			//2.output 에 creatorTxs 안에 있는 트랜잭션 내에 없다는 거을 확인
			for index, output := range tx.TxOuts {
				//ok는 tx.ID key 값이 map 안에 있는지 없는지 확인해주는 값
				//if output.Owner == address {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						//추가하려는 unspent transaction output 이 mempool에 아직 없는지 확인
						uTxOut := &UTxOut{
							TxID:   tx.ID,
							Index:  index,
							Amount: output.Amount,
						}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

//3.단하나의 주소만 보여주도록

//func (b blockchain) BalanceByAddress(address string) int {
func BalanceByAddress(b *blockchain, address string) int {
	//txOuts := b.TxOutsByAddress(address)
	txOuts := UTxOutsByAddress(b, address)
	var result int
	for _, tx := range txOuts {
		result += tx.Amount
	}
	return result
}
