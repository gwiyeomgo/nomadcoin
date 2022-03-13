package blockchain

import (
	"github.com/gwiyeomgo/nomadcoin/db"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newest_hash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

//블록이 새로 만들어 질떄마다 블록과 블록체인의 상황을 db에 저장
//db에 block 을 저장하는 코드
//블록체인을 처음 실행시키는 사람의 관점
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height
}

func Blockchain() *blockchain {
	if b == nil {
		//b를 초기화 (처음이자 마지막)
		once.Do(func() {
			//우리가 몇 천 개의 goroutine을 실행해도 오직 한번만 호출
			b = &blockchain{"", 0}
			b.AddBlock("Genesis")
		})
	}
	return b
}
