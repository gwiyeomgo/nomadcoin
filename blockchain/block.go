package blockchain

import (
	"crypto/sha256"
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/db"
	"github.com/gwiyeomgo/nomadcoin/utils"
)

type Block struct {
	Data    string `json:"data"`
	Hash    string `json:"hash"`
	PreHash string `json:"prehash,omitempty"`
	Height  int    `json:"height"`
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

func createBlock(data string, preHash string, height int) Block {
	block := Block{
		Data:    data,
		Hash:    "",
		PreHash: preHash,
		Height:  height,
	}
	payload := block.Data + block.PreHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	//(1) persist 는 SaveBlock 을 호출
	block.persist()
	return block
}
