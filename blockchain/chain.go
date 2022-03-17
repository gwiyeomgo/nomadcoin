package blockchain

import (
	"fmt"
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

//block에 checkoutpoint 를 지정
func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

//블록이 새로 만들어 질떄마다 블록과 블록체인의 상황을 db에 저장
//db에 block 을 저장하는 코드
//블록체인을 처음 실행시키는 사람의 관점
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		//b를 초기화 (처음이자 마지막)
		once.Do(func() {
			//우리가 몇 천 개의 goroutine을 실행해도 오직 한번만 호출

			b = &blockchain{"", 0}
			// search for checkpoint on the db (db에서 블록체인을 가져오는 함수)
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Println("...Restore")
				//restore b from bytes (db 에는 bytes로 저장되어 있음)
				b.restore(checkpoint)
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}

//모든 block 을 보여주는 함수
//Block 포인터의 slice 를 retrun
func (b *blockchain) Blocks() []*Block {
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
