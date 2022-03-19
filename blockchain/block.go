package blockchain

import (
	"errors"
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/db"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"strings"
	"time"
)

//const difficulty int = 2
//Block 에 Data는 지우고
// Transaction 추가
type Block struct {
	//	Data       string `json:"data"`
	Hash        string `json:"hash"`
	PreHash     string `json:"prehash,omitempty"`
	Height      int    `json:"height"`
	Difficulty  int    `json:"difficulty"`
	Nonce       int    `json:"nonce"`
	Timestamp   int    `json:"timestamp"` // 블록 당 생성시간을 알 수 있음
	Transaction []*Tx  `json:"transaction"`
}

//utils로 빼줌
/*func (b *Block) toBytes() []byte {
	//golang 의 value 를 byte encode 나 decode 하는 패키지 gob
	// encoder 를 만들고 block 을 encode 한다음
	// 그 결과를 blockBuffer 에 저장
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	utils.HandleErr(encoder.Encode(b))
	return blockBuffer.Bytes()
}*/

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrNotFound = errors.New("Block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

//func createBlock(data string, preHash string, height int) Block {
func createBlock(preHash string, height int) Block {
	block := Block{
		//		Data:       data,
		Hash:        "",
		PreHash:     preHash,
		Height:      height,
		Difficulty:  Blockchain().difficulty(),
		Nonce:       0,
		Transaction: []*Tx{makeCoinbaseTx("gwiyeom")},
	}
	/*	payload := block.Data + block.PreHash + fmt.Sprint(block.Height)
		block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	*/
	block.mine()
	//(1) persist 는 SaveBlock 을 호출
	block.persist()
	return block
}
func (b *Block) mine() {
	//해쉬의 시작에 0이 몇 개 있는지 확인
	// string 을 difficultly 번 반복 => 00 출력
	target := strings.Repeat("0", b.Difficulty)
	//nonce 를 개속 늘려 원하는 값을 찾는다.
	for {
		//블록을 hash 하기 전에 블록 타임스템프를 지정하고
		b.Timestamp = int(time.Now().Unix()) //.Unix() 는 int64 를 반환
		//hash 를 16진수의 string 으로 변경
		//blockAsString := fmt.Sprint(b)
		//hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))
		hash := utils.Hash(b)
		fmt.Printf("Target:%s\nHash:%s\nnonce:%s\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}
