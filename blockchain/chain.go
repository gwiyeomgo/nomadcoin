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

//block에 checkoutpoint 를 지정
func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

//블록이 새로 만들어 질떄마다 블록과 블록체인의 상황을 db에 저장
//db에 block 을 저장하는 코드
//블록체인을 처음 실행시키는 사람의 관점
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	//difficulty 지정
	//block 추가할때마다 chain에 difficulty 변경됨
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

//예상 시간은 5개의 블록이 2분마다 생성되는 시간 => 10분
//10분동안 5개 블록 생성되길 예상
func (b *blockchain) recalculateDifficulty() int {
	//모든 blocks 를 받아온다
	allBlocks := b.Blocks()
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
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//만약 나머지값이 0이라면
		// recalculate the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func Blockchain() *blockchain {
	if b == nil {
		//b를 초기화 (처음이자 마지막)
		once.Do(func() {
			//우리가 몇 천 개의 goroutine을 실행해도 오직 한번만 호출

			b = &blockchain{
				Height: 0,
			}
			// search for checkpoint on the db (db에서 블록체인을 가져오는 함수)
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				//fmt.Println("...Restore")
				//restore b from bytes (db 에는 bytes로 저장되어 있음)
				b.restore(checkpoint)
			}
		})
	}
	//	fmt.Printf("NewesHash:%s\nHeight:%d\n", b.NewestHash, b.Height)
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
